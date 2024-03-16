package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	// Import godotenv package
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get database connection string from environment variable
	connectionString := os.Getenv("DB_CONNECTION_URL")
	if connectionString == "" {
		connectionString = "postgresql://psql-db_owner:l1C0IkXSrKqz@ep-fragrant-glade-a2d997l1.eu-central-1.aws.neon.tech/psql-db?sslmode=require"
	}

	// Initialize the database connection
	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verify the database connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", homeHandler)
	router.Get("/jobs", jobsHandler)

	log.Println("Server started on port 8080")
	http.ListenAndServe(":"+port, router)
}

// Rest of your handlers remain the same...

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

func jobsHandler(w http.ResponseWriter, r *http.Request) {
	// Query jobs from the database
	rows, err := db.Query("SELECT * FROM jobs")
	if err != nil {
		http.Error(w, "Failed to query jobs", http.StatusInternalServerError)
		log.Println("Failed to query jobs:", err)
		return
	}
	defer rows.Close()

	// Define a struct to represent a job
	type Job struct {
		JobID       int     `json:"job_id"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		Company     string  `json:"company"`
		Salary      float64 `json:"salary"`
		PostedDate  string  `json:"posted_date"`
		Deadline    string  `json:"deadline_date"`
	}

	// Create a slice to hold the job data
	var jobs []Job

	// Iterate over the rows and populate the job slice
	for rows.Next() {
		var job Job
		err := rows.Scan(&job.JobID, &job.Title, &job.Description, &job.Latitude, &job.Longitude, &job.Company, &job.Salary, &job.PostedDate, &job.Deadline)
		if err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			log.Println("Failed to scan row:", err)
			return
		}
		jobs = append(jobs, job)
	}

	// Marshal the job slice into JSON format
	jsonData, err := json.Marshal(jobs)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		log.Println("Failed to marshal JSON:", err)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON data to the response
	w.Write(jsonData)
}
