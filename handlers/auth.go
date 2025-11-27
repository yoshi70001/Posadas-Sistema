package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"

	"posadas-sistema/database"
	"posadas-sistema/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("mi_clave_secreta_super_larga_y_compleja_para_jwt")

// Claims is a struct that will be encoded to a JWT.
// We add `jwtRegisteredClaims` which is a jwt.RegisteredClaims that contains standard claims.
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var user models.User
		var isActive bool
		err := database.DB.QueryRow("SELECT id, username, password, is_active FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Password, &isActive)
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		user.IsActive = isActive
		if !user.IsActive {
			http.Error(w, "Account is inactive", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Create JWT Token
		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &Claims{
			UserID:   user.ID,
			Username: user.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "jwt_token",
			Value:   tokenString,
			Expires: expirationTime,
		})

		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}
	// Render login template
	tmpl, err := template.ParseFiles("templates/base.html", "templates/login.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt_token",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("jwt_token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenString := c.Value
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if !token.Valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Verify user is still active in the database
		var isActive bool
		err = database.DB.QueryRow("SELECT is_active FROM users WHERE id = ?", claims.UserID).Scan(&isActive)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		if !isActive {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next(w, r)
	}
}
