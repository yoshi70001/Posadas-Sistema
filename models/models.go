package models

import "time"

type Registration struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Age             int       `json:"age"`
	DNI             string    `json:"dni"`
	GuardianName    string    `json:"guardian_name"`
	GuardianContact string    `json:"guardian_contact"`
	Year            int       `json:"year"`
	CreatedAt       time.Time `json:"created_at"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // Hashed
	IsActive bool   `json:"is_active"`
}

type Event struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"` // "ensayo" o "salida"
	Date        time.Time `json:"date"`
	Time        string    `json:"time"` // "4:00 PM" o "6:00 PM"
	Location    string    `json:"location"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Attendance struct {
	ID            int       `json:"id"`
	EventID       int       `json:"event_id"`
	RegistrationID int      `json:"registration_id"`
	Present       bool      `json:"present"`
	Notes         string    `json:"notes"`
	MarkedAt      time.Time `json:"marked_at"`
}
