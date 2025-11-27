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
