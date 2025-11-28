package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"posadas-sistema/database"
	"posadas-sistema/models"

	"golang.org/x/crypto/bcrypt"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, age, dni, guardian_name, guardian_contact, year, created_at FROM registrations ORDER BY created_at DESC")
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var registrations []models.Registration
	for rows.Next() {
		var reg models.Registration
		if err := rows.Scan(&reg.ID, &reg.Name, &reg.Age, &reg.DNI, &reg.GuardianName, &reg.GuardianContact, &reg.Year, &reg.CreatedAt); err != nil {
			log.Println(err)
			continue
		}
		registrations = append(registrations, reg)
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/dashboard.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Registrations []models.Registration
		User          string
	}{
		Registrations: registrations,
		User:          "Admin", // You could get this from the session
	}

	tmpl.Execute(w, data)
}

// AdminListHandler lists all admin users
func AdminListHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, username, is_active FROM users")
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		var isActive bool
		if err := rows.Scan(&u.ID, &u.Username, &isActive); err != nil {
			log.Println(err)
			continue
		}
		u.IsActive = isActive
		users = append(users, u)
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/admin_list.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, users)
}

// AdminCreateHandler shows the form to create a new admin
func AdminCreateHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/admin_form.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// AdminStoreHandler saves the new admin
func AdminStoreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec("INSERT INTO users (username, password, is_active) VALUES (?, ?, ?)", username, hashedPassword, true)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

// ===== EVENT MANAGEMENT HANDLERS =====

// EventListHandler lists all events
func EventListHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, type, date, time, location, description, created_at FROM events ORDER BY date DESC")
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.Name, &event.Type, &event.Date, &event.Time, &event.Location, &event.Description, &event.CreatedAt); err != nil {
			log.Println(err)
			continue
		}
		events = append(events, event)
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/events_list.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, events)
}

// EventCreateHandler shows the form to create a new event
func EventCreateHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/events_form.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// EventStoreHandler saves the new event
func EventStoreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	eventType := r.FormValue("type")
	date := r.FormValue("date")
	time := r.FormValue("time")
	location := r.FormValue("location")
	description := r.FormValue("description")

	_, err := database.DB.Exec("INSERT INTO events (name, type, date, time, location, description) VALUES (?, ?, ?, ?, ?, ?)",
		name, eventType, date, time, location, description)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/events", http.StatusSeeOther)
}

// EventEditHandler shows the form to edit an event
func EventEditHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var event models.Event
	err := database.DB.QueryRow("SELECT id, name, type, date, time, location, description FROM events WHERE id = ?", id).
		Scan(&event.ID, &event.Name, &event.Type, &event.Date, &event.Time, &event.Location, &event.Description)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/events_form.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, event)
}

// EventUpdateHandler updates the event
func EventUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	name := r.FormValue("name")
	eventType := r.FormValue("type")
	date := r.FormValue("date")
	time := r.FormValue("time")
	location := r.FormValue("location")
	description := r.FormValue("description")

	_, err := database.DB.Exec("UPDATE events SET name = ?, type = ?, date = ?, time = ?, location = ?, description = ? WHERE id = ?",
		name, eventType, date, time, location, description, id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/events", http.StatusSeeOther)
}

// EventDeleteHandler deletes an event
func EventDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := database.DB.Exec("DELETE FROM events WHERE id = ?", id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/events", http.StatusSeeOther)
}

// ===== ATTENDANCE HANDLERS =====

// AttendanceHandler shows the form to mark attendance for an event
func AttendanceHandler(w http.ResponseWriter, r *http.Request) {
	eventID := r.URL.Query().Get("event_id")
	if eventID == "" {
		http.Error(w, "Event ID required", http.StatusBadRequest)
		return
	}

	// Get event details
	var event models.Event
	err := database.DB.QueryRow("SELECT id, name, type, date, time, location FROM events WHERE id = ?", eventID).
		Scan(&event.ID, &event.Name, &event.Type, &event.Date, &event.Time, &event.Location)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	// Get all registrations
	rows, err := database.DB.Query("SELECT id, name, age, dni FROM registrations ORDER BY name")
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type AttendanceForm struct {
		Event         models.Event
		Registrations []struct {
			ID              int    `json:"id"`
			Name            string `json:"name"`
			Age             int    `json:"age"`
			DNI             string `json:"dni"`
			Present         bool   `json:"present"`
			HasAttendance   bool   `json:"has_attendance"`
		}
	}

	var formData AttendanceForm
	formData.Event = event

	for rows.Next() {
		var reg struct {
			ID              int    `json:"id"`
			Name            string `json:"name"`
			Age             int    `json:"age"`
			DNI             string `json:"dni"`
			Present         bool   `json:"present"`
			HasAttendance   bool   `json:"has_attendance"`
		}

		if err := rows.Scan(&reg.ID, &reg.Name, &reg.Age, &reg.DNI); err != nil {
			log.Println(err)
			continue
		}

		// Check if attendance already exists for this registration and event
		var present bool
		err := database.DB.QueryRow("SELECT present FROM attendance WHERE event_id = ? AND registration_id = ?", eventID, reg.ID).
			Scan(&present)
		if err == nil {
			reg.Present = present
			reg.HasAttendance = true
		} else {
			reg.Present = false
			reg.HasAttendance = false
		}

		formData.Registrations = append(formData.Registrations, reg)
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/attendance_form.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, formData)
}

// AttendanceStoreHandler saves the attendance data
func AttendanceStoreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	eventID := r.FormValue("event_id")
	presentRegistrations := r.Form["present"] // This gets all values for "present" field
	notes := r.FormValue("notes")

	// First, get all current registrations to know who was present
	var allRegistrations []int
	rows, err := database.DB.Query("SELECT id FROM registrations")
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var regID int
		if err := rows.Scan(&regID); err != nil {
			log.Println(err)
			continue
		}
		allRegistrations = append(allRegistrations, regID)
	}

	// Create a map of present registrations
	presentMap := make(map[int]bool)
	for _, regStr := range presentRegistrations {
		var regID int
		if _, err := fmt.Sscanf(regStr, "%d", &regID); err == nil {
			presentMap[regID] = true
		}
	}

	// Delete existing attendance for this event (to replace with new data)
	_, err = database.DB.Exec("DELETE FROM attendance WHERE event_id = ?", eventID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Insert attendance records
	for _, regID := range allRegistrations {
		present := presentMap[regID]
		_, err := database.DB.Exec("INSERT INTO attendance (event_id, registration_id, present, notes) VALUES (?, ?, ?, ?)",
			eventID, regID, present, notes)
		if err != nil {
			log.Println(err)
			// Continue with other registrations even if one fails
		}
	}

	http.Redirect(w, r, "/admin/events", http.StatusSeeOther)
}

// DashboardDataHandler returns JSON data for the dashboard
func DashboardDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get attendance statistics
	var totalEvents int
	var totalAttendances int
	var ensayosCount int
	var salidasCount int

	database.DB.QueryRow("SELECT COUNT(*) FROM events").Scan(&totalEvents)
	database.DB.QueryRow("SELECT COUNT(*) FROM attendance WHERE present = 1").Scan(&totalAttendances)
	database.DB.QueryRow("SELECT COUNT(*) FROM events WHERE type = 'ensayo'").Scan(&ensayosCount)
	database.DB.QueryRow("SELECT COUNT(*) FROM events WHERE type = 'salida'").Scan(&salidasCount)

	// Get monthly attendance data
	monthlyQuery := `
		SELECT
			strftime('%Y-%m', e.date) as month,
			e.type,
			COUNT(CASE WHEN a.present = 1 THEN 1 END) as attendances,
			COUNT(*) as total_registrations
		FROM events e
		LEFT JOIN attendance a ON e.id = a.event_id
		GROUP BY strftime('%Y-%m', e.date), e.type
		ORDER BY month`

	rows, err := database.DB.Query(monthlyQuery)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type MonthlyData struct {
		Month           string `json:"month"`
		Type            string `json:"type"`
		Attendances     int    `json:"attendances"`
		TotalRegistrations int `json:"total_registrations"`
	}

	var monthlyData []MonthlyData
	for rows.Next() {
		var data MonthlyData
		if err := rows.Scan(&data.Month, &data.Type, &data.Attendances, &data.TotalRegistrations); err != nil {
			log.Println(err)
			continue
		}
		monthlyData = append(monthlyData, data)
	}

	// Calculate percentage for monthly data
	for i := range monthlyData {
		if monthlyData[i].TotalRegistrations > 0 {
			// This would need to be calculated properly - for now return raw data
		}
	}

	response := struct {
		TotalEvents        int             `json:"total_events"`
		TotalAttendances   int             `json:"total_attendances"`
		EnsayosCount       int             `json:"ensayos_count"`
		SalidasCount       int             `json:"salidas_count"`
		MonthlyData        []MonthlyData   `json:"monthly_data"`
	}{
		TotalEvents:      totalEvents,
		TotalAttendances: totalAttendances,
		EnsayosCount:     ensayosCount,
		SalidasCount:     salidasCount,
		MonthlyData:      monthlyData,
	}

	json.NewEncoder(w).Encode(response)
}

// AdminEditHandler shows the form to edit an admin
func AdminEditHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var user models.User
	var isActive bool
	err := database.DB.QueryRow("SELECT id, username, is_active FROM users WHERE id = ?", id).Scan(&user.ID, &user.Username, &isActive)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	user.IsActive = isActive

	tmpl, err := template.ParseFiles("templates/base.html", "templates/admin_form.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, user)
}

// AdminUpdateHandler updates the admin
func AdminUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	username := r.FormValue("username")
	password := r.FormValue("password")

	var err error
	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		_, err = database.DB.Exec("UPDATE users SET username = ?, password = ? WHERE id = ?", username, hashedPassword, id)
	} else {
		_, err = database.DB.Exec("UPDATE users SET username = ? WHERE id = ?", username, id)
	}

	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

// AdminDeleteHandler deletes an admin
func AdminDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := database.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

// AdminToggleStatusHandler toggles the active status of an admin
func AdminToggleStatusHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var currentStatus bool
	err := database.DB.QueryRow("SELECT is_active FROM users WHERE id = ?", id).Scan(&currentStatus)
	if err != nil {
		log.Println(err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	newStatus := !currentStatus // Toggle the status

	_, err = database.DB.Exec("UPDATE users SET is_active = ? WHERE id = ?", newStatus, id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
