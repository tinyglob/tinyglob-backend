package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"

	handlers "tinyglob-backend/controllers"
	"tinyglob-backend/database"
	middlewares "tinyglob-backend/middleware" // Import your custom middleware package
)

func main() {
	db_instance := database.InitDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Middleware
	corsMiddleware := cors.Default()
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(corsMiddleware.Handler)

	// Apply AuthMiddleware to routes that require authentication
	router.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)
		// Protected routes
		r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("This is a protected route\n"))
		})
	})

	// Generic
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to TinyGlob!\n"))
	})

	// Jobs
	router.Get("/jobs", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetJobsCount(db_instance, w, r)
	})

	router.Get("/jobs/continent/{continent}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetJobsByContinent(db_instance, w, r)
	})

	router.Get("/jobs/country/{country}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetJobsByCountry(db_instance, w, r)
	})

	router.Get("/jobs/id/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetJobById(db_instance, w, r)
	})

	// User
	router.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.Register(db_instance, w, r)
	})

	router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.Login(db_instance, w, r)
	})

	router.Post("/token", func(w http.ResponseWriter, r *http.Request) {
		handlers.GenerateToken(db_instance, w, r)
	})

	log.Printf("Server is started on port %s", port)
	http.ListenAndServe(":"+port, router)
}
