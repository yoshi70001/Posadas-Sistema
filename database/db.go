package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./posadas.db")
	if err != nil {
		log.Fatal(err)
	}

	createTables()
}

func createTables() {
	createRegistrationsTable := `
	CREATE TABLE IF NOT EXISTS registrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age INTEGER NOT NULL,
		dni TEXT NOT NULL,
		guardian_name TEXT,
		guardian_contact TEXT,
		year INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		is_active BOOLEAN NOT NULL
	);`

	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		type TEXT NOT NULL CHECK(type IN ('ensayo', 'salida')),
		date DATE NOT NULL,
		time TEXT NOT NULL,
		location TEXT NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createAttendanceTable := `
	CREATE TABLE IF NOT EXISTS attendance (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_id INTEGER NOT NULL,
		registration_id INTEGER NOT NULL,
		present BOOLEAN NOT NULL DEFAULT 0,
		notes TEXT,
		marked_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(event_id) REFERENCES events(id) ON DELETE CASCADE,
		FOREIGN KEY(registration_id) REFERENCES registrations(id) ON DELETE CASCADE,
		UNIQUE(event_id, registration_id)
	);`

	_, err := DB.Exec(createRegistrationsTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(createUsersTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(createEventsTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(createAttendanceTable)
	if err != nil {
		log.Fatal(err)
	}

	seedAdmin()
}

func seedAdmin() {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Println("Error checking users:", err)
		return
	}

	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error hashing admin password:", err)
			return
		}

		_, err = DB.Exec("INSERT INTO users (username, password, is_active) VALUES (?, ?, ?)", "admin", hashedPassword, true)
		if err != nil {
			log.Println("Error seeding admin:", err)
		} else {
			log.Println("Default admin user created (admin/admin)")
		}
	}
}
