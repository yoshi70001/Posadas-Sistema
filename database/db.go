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

	_, err := DB.Exec(createRegistrationsTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(createUsersTable)
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
