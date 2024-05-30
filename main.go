package main

import (
	"log"
	"net/http"
	"os"
	"tinyglob-backend/controllers"
	"tinyglob-backend/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func main() {

	db_instance := db.InitDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	corsMiddleware := cors.Default()
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(corsMiddleware.Handler)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to TinyGlob!\n"))
	})

	router.Get("/jobs", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetJobsCount(db_instance, w, r)
	})

	router.Get("/jobs/continent/{continent}", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetJobsByContinent(db_instance, w, r)
	})

	router.Get("/jobs/country/{country}", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetJobsByCountry(db_instance, w, r)
	})

	router.Get("/jobs/id/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetJobById(db_instance, w, r)
	})

	log.Printf("Server is started on port %s", port)
	http.ListenAndServe(":"+port, router)
}
