package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/rs/cors" // Import rs/cors for CORS handling
)

var db *sql.DB

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	connectionString := os.Getenv("DB_CONNECTION_URL")

	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// Create a CORS middleware instance with default options
	corsMiddleware := cors.Default()

	// Use CORS middleware with your router
	router.Use(corsMiddleware.Handler)

	router.Get("/", homeHandler)
	router.Get("/jobs", jobsHandler)

	log.Println("Server started on port 8080")
	http.ListenAndServe(":"+port, router)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

func jobsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM jobs")
	if err != nil {
		http.Error(w, "Failed to query jobs", http.StatusInternalServerError)
		log.Println("Failed to query jobs:", err)
		return
	}
	defer rows.Close()

	type Job struct {
		JobID          int      `json:"job_id"`
		VideoID        int      `json:"video_id"`
		VideoUrl       int      `json:"video_url"`
		Title          string   `json:"title"`
		Description    string   `json:"description"`
		Continent      string   `json:"continent"`
		Country        string   `json:"country"`
		City           string   `json:"city"`
		Company        string   `json:"company"`
		CompanyLogoUrl string   `json:"company_logo_url"`
		Salary         float64  `json:"salary"`
		Currency       string   `json:"currency"`
		RequiredSkills []string `json:"required_skills"`
		PostedDate     string   `json:"posted_date"`
		Deadline       string   `json:"deadline_date"`
	}

	var jobs []Job

	for rows.Next() {
		var job Job
		var requiredSkillsString string
		err := rows.Scan(&job.JobID, &job.VideoID, &job.VideoUrl, &job.Title, &job.Description, &job.Continent, &job.Country, &job.City, &job.Company, &job.CompanyLogoUrl, &job.Salary, &job.Currency, &requiredSkillsString, &job.PostedDate, &job.Deadline)
		if err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			log.Println("Failed to scan row:", err)
			return
		}
		job.RequiredSkills = strings.Split(requiredSkillsString, ",") // Split the string into an array of strings
		jobs = append(jobs, job)
	}

	jsonData, err := json.Marshal(jobs)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		log.Println("Failed to marshal JSON:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
