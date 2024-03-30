package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"tinyglob-backend/types"
)

// curl -X POST http://localhost:8080/register -d '{"name": "Nick", "username": "nickydicky", "email": "nickydicky@gmail.com", "password": "123"}'

func Register(db_instance *sql.DB, w http.ResponseWriter, r *http.Request) {
	var user types.User

	// Decode the JSON body into the user struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if all required fields are present
	if user.Name == "" || user.Username == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "All parameters (name, username, email, password) are required", http.StatusBadRequest)
		return
	}
	// Check if email format is valid
	if !isValidEmail(user.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Hash the password
	err = user.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Failed to hash the password", http.StatusInternalServerError)
		return
	}

	// Insert the user into the database
	_, err = db_instance.Exec("INSERT INTO users (name, username, email, hashed_password) VALUES ($1, $2, $3, $4)",
		user.Name, user.Username, user.Email, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// You can now use the user object
	fmt.Println("User registered successfully:", user)

	// Optionally, you can return a success message
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User registered successfully")
}

func Login(db_instance *sql.DB, w http.ResponseWriter, r *http.Request) {
	var request types.TokenRequest

	// Decode the JSON body into the TokenRequest struct
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Query the database to retrieve the user with the provided email
	var user types.User
	err = db_instance.QueryRow("SELECT name, username, email, hashed_password FROM users WHERE email = $1", request.Email).Scan(&user.Name, &user.Username, &user.Email, &user.Password)
	if err != nil {
		log.Println("Error querying database:", err)
		http.Error(w, "Email not found", http.StatusUnauthorized)
		return
	}

	// Check if the provided password is correct
	err = user.CheckPassword(request.Password)
	if err != nil {
		http.Error(w, "Password is incorrect", http.StatusUnauthorized)
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
