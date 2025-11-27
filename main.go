package main

import (
	"log"
	"net/http"

	"posadas-sistema/database"
	"posadas-sistema/handlers"
)

func main() {
	database.InitDB()

	// Static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Public Routes
	http.HandleFunc("/", handlers.LandingHandler)
	http.HandleFunc("/register", handlers.RegisterFormHandler)
	http.HandleFunc("/register/submit", handlers.RegisterSubmitHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	// Admin Routes (Protected)
	http.HandleFunc("/admin/dashboard", handlers.AuthMiddleware(handlers.DashboardHandler))
	http.HandleFunc("/admin/users", handlers.AuthMiddleware(handlers.AdminListHandler))
	http.HandleFunc("/admin/users/create", handlers.AuthMiddleware(handlers.AdminCreateHandler))
	http.HandleFunc("/admin/users/store", handlers.AuthMiddleware(handlers.AdminStoreHandler))
	http.HandleFunc("/admin/users/edit", handlers.AuthMiddleware(handlers.AdminEditHandler))
	http.HandleFunc("/admin/users/update", handlers.AuthMiddleware(handlers.AdminUpdateHandler))
	http.HandleFunc("/admin/users/delete", handlers.AuthMiddleware(handlers.AdminDeleteHandler))
	http.HandleFunc("/admin/users/toggle-status", handlers.AuthMiddleware(handlers.AdminToggleStatusHandler))

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
