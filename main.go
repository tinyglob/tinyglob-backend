package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

type Job struct {
	JobID          int      `json:"job_id"`
	VideoID        int      `json:"video_id"`
	VideoUrl       string   `json:"video_url"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Country        string   `json:"country"`
	City           string   `json:"city"`
	Continent      string   `json:"continent"`
	Company        string   `json:"company"`
	CompanyLogoUrl string   `json:"company_logo_url"`
	Salary         float64  `json:"salary"`
	Currency       string   `json:"currency"`
	RequiredSkills []string `json:"required_skills"`
	PostedDate     string   `json:"posted_date"`
	Deadline       string   `json:"deadline_date"`
}

var db *sql.DB

func initDB() {
	connectionString := os.Getenv("DB_CONNECTION_URL")
	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(10) // Set maximum open connections
	db.SetMaxIdleConns(5)  // Set maximum idle connections
}

func main() {
	initDB() // Initialize the database connection pool

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
	router.Get("/", getRootHandler)
	router.Get("/jobs", getAllJobsHandler)
	router.Get("/jobs/id/{id}", getJobByIDHandler)                       // Specify route pattern with "id" parameter
	router.Get("/jobs/continent/{continent}", getJobsByContinentHandler) // Specify route pattern with "continent" parameter

	log.Printf("Server started on port %s", port)
	http.ListenAndServe(":"+port, router)
}

func getRootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

func getAllJobsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT job_id, video_id, video_url, title, continent, country, city, company, company_logo_url FROM jobs")
	if err != nil {
		http.Error(w, "Failed to query jobs", http.StatusInternalServerError)
		log.Println("Failed to query jobs:", err)
		return
	}
	defer rows.Close()

	var jobs []Job

	for rows.Next() {
		var job Job
		err := rows.Scan(&job.JobID, &job.VideoID, &job.VideoUrl, &job.Title, &job.Continent, &job.Country, &job.City, &job.Company, &job.CompanyLogoUrl)
		if err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			log.Println("Failed to scan row:", err)
			return
		}
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

func getJobByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Get the job ID from the URL parameters
	jobID := chi.URLParam(r, "id")
	log.Println("Received request for job ID:", jobID)

	// Convert jobID to int
	jobIDInt, err := strconv.Atoi(jobID)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		log.Println("Invalid job ID:", err)
		return
	}

	// Query the database to retrieve the job with the given ID
	log.Println("Querying database for job with ID:", jobIDInt)
	row := db.QueryRow("SELECT job_id, video_id, video_url, title, description, continent, country, city, company, salary, currency, posted_date, deadline_date FROM jobs WHERE job_id = $1", jobIDInt)

	var job Job
	err = row.Scan(&job.JobID, &job.VideoID, &job.VideoUrl, &job.Title, &job.Description, &job.Continent, &job.Country, &job.City, &job.Company, &job.Salary, &job.Currency, &job.PostedDate, &job.Deadline)
	if err != nil {
		http.Error(w, "Failed to retrieve job", http.StatusInternalServerError)
		log.Println("Failed to retrieve job:", err)
		return
	}

	// Convert into JSON format
	jsonData, err := json.Marshal(job)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		log.Println("Failed to marshal JSON:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func getJobsByContinentHandler(w http.ResponseWriter, r *http.Request) {
	continent := chi.URLParam(r, "continent")
	if continent == "" {
		http.Error(w, "Continent parameter is required", http.StatusBadRequest)
		return
	}

	// Convert continent to lowercase
	continent = strings.ToLower(continent)

	rows, err := db.Query("SELECT job_id, video_id, video_url, title, description, continent, country, city, company, company_logo_url, salary, currency FROM jobs WHERE LOWER(continent) = $1", continent)
	if err != nil {
		http.Error(w, "Failed to query jobs", http.StatusInternalServerError)
		log.Println("Failed to query jobs:", err)
		return
	}
	defer rows.Close()

	var jobs []Job

	for rows.Next() {
		var job Job
		err := rows.Scan(&job.JobID, &job.VideoID, &job.VideoUrl, &job.Title, &job.Description, &job.Continent, &job.Country, &job.City, &job.Company, &job.CompanyLogoUrl, &job.Salary, &job.Currency)
		if err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			log.Println("Failed to scan row:", err)
			return
		}
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
