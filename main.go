package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"

	"tinyglob-backend/database"
	"tinyglob-backend/handlers"
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

	// Routes
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to TinyGlob!\n"))
	})

	router.Get("/jobs", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetJobsCount(db_instance, w, r)
	})

	router.Get("/jobs/id/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetJobById(db_instance, w, r)
	})

	router.Get("/jobs/continent/{continent}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetJobsByContinent(db_instance, w, r)
	})

	log.Printf("Server is started on port %s", port)
	http.ListenAndServe(":"+port, router)
}
