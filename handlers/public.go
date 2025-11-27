package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"posadas-sistema/database"
)

func LandingHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/index.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func RegisterFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/register.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func RegisterSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	ageStr := r.FormValue("age")
	dni := r.FormValue("dni")
	guardianName := r.FormValue("guardian_name")
	guardianContact := r.FormValue("guardian_contact")
	yearStr := r.FormValue("year")

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		http.Error(w, "Invalid age", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	stmt, err := database.DB.Prepare("INSERT INTO registrations(name, age, dni, guardian_name, guardian_contact, year) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, age, dni, guardianName, guardianContact, year)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Redirect to success page or back to home with message
	http.Redirect(w, r, "/?success=true", http.StatusSeeOther)
}
