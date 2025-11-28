package main

import (
	"log"
	"net/http"

	"posadas-sistema/database"
	"posadas-sistema/handlers"
)

func main() {
	database.InitDB()

	mux := http.NewServeMux()

	// Static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// Public Routes
	mux.HandleFunc("GET /", handlers.LandingHandler)
	mux.HandleFunc("GET /register", handlers.RegisterFormHandler)
	mux.HandleFunc("POST /register/submit", handlers.RegisterSubmitHandler)
	mux.HandleFunc("GET /login", handlers.LoginHandler)
	mux.HandleFunc("POST /login", handlers.LoginHandler) // Allow POST for login submission
	mux.HandleFunc("GET /logout", handlers.LogoutHandler)

	// Admin Routes (Protected)
	mux.HandleFunc("GET /admin/dashboard", handlers.AuthMiddleware(handlers.DashboardHandler))
	mux.HandleFunc("GET /admin/dashboard-data", handlers.AuthMiddleware(handlers.DashboardDataHandler))
	mux.HandleFunc("GET /admin/users", handlers.AuthMiddleware(handlers.AdminListHandler))
	mux.HandleFunc("GET /admin/users/create", handlers.AuthMiddleware(handlers.AdminCreateHandler))
	mux.HandleFunc("POST /admin/users/store", handlers.AuthMiddleware(handlers.AdminStoreHandler))
	mux.HandleFunc("GET /admin/users/edit", handlers.AuthMiddleware(handlers.AdminEditHandler))
	mux.HandleFunc("POST /admin/users/update", handlers.AuthMiddleware(handlers.AdminUpdateHandler))
	mux.HandleFunc("POST /admin/users/delete", handlers.AuthMiddleware(handlers.AdminDeleteHandler))
	mux.HandleFunc("GET /admin/users/toggle-status", handlers.AuthMiddleware(handlers.AdminToggleStatusHandler))

	// Event Management Routes (Protected)
	mux.HandleFunc("GET /admin/events", handlers.AuthMiddleware(handlers.EventListHandler))
	mux.HandleFunc("GET /admin/events/create", handlers.AuthMiddleware(handlers.EventCreateHandler))
	mux.HandleFunc("POST /admin/events/store", handlers.AuthMiddleware(handlers.EventStoreHandler))
	mux.HandleFunc("GET /admin/events/edit", handlers.AuthMiddleware(handlers.EventEditHandler))
	mux.HandleFunc("POST /admin/events/update", handlers.AuthMiddleware(handlers.EventUpdateHandler))
	mux.HandleFunc("POST /admin/events/delete", handlers.AuthMiddleware(handlers.EventDeleteHandler))

	// Attendance Routes (Protected)
	mux.HandleFunc("GET /admin/attendance", handlers.AuthMiddleware(handlers.AttendanceHandler))
	mux.HandleFunc("POST /admin/attendance/store", handlers.AuthMiddleware(handlers.AttendanceStoreHandler))

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
