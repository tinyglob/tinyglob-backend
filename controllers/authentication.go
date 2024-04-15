package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"tinyglob-backend/types"
)

func GenerateToken(db_instance *sql.DB, w http.ResponseWriter, r *http.Request) {
	var request types.TokenRequest
	var user types.User

	// Decode the JSON body into the TokenRequest struct
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate email format
	if !isValidEmail(request.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Check if email exists in the database
	row := db_instance.QueryRow("SELECT name, username, email, hashed_password FROM users WHERE email = $1", request.Email)
	if err := row.Scan(&user.Name, &user.Username, &user.Email, &user.Password); err != nil {
		http.Error(w, "Email not found: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Check if the provided password is correct
	if err := user.CheckPassword(request.Password); err != nil {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	tokenString, err := types.GenerateJWT(user.Email, user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Respond with the token
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tokenString))
}

// isValidEmail checks if the email format is valid
func isValidEmail(email string) bool {
	// Basic email validation using regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
