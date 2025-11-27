package handlers

import (
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
