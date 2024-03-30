package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"tinyglob-backend/helpers"
	"tinyglob-backend/types"
)

func GetJobById(db_instance *sql.DB, w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "id")

	jobIDInt, err := strconv.Atoi(jobID)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		log.Println("Invalid job ID:", err)
		return
	}

	stmt, err := db_instance.Prepare("SELECT job_id, video_id, video_url, title, description, continent, country, category, city, company, start_salary, end_salary, currency, posted_date, deadline_date FROM jobs WHERE job_id = $1")
	if err != nil {
		http.Error(w, "Failed to prepare SQL statement", http.StatusInternalServerError)
		log.Println("Failed to prepare SQL statement:", err)
		return
	}
	defer stmt.Close()

	row := stmt.QueryRow(jobIDInt)

	var job types.Job
	err = row.Scan(&job.JobID, &job.VideoID, &job.VideoUrl, &job.Title, &job.Description, &job.Country, &job.Category, &job.City, &job.Continent, &job.Company, &job.CompanyLogoUrl, &job.StartSalary, &job.EndSalary, &job.Currency, &job.PostedDate, &job.Deadline)
	if err != nil {
		http.Error(w, "Failed to retrieve job", http.StatusInternalServerError)
		log.Println("Failed to retrieve job:", err)
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, job)
}

func GetJobsCount(db_instance *sql.DB, w http.ResponseWriter, r *http.Request) {
	rows, err := db_instance.Query("SELECT continent, COUNT(*) AS job_count FROM jobs WHERE continent IS NOT NULL GROUP BY continent")
	if err != nil {
		http.Error(w, "Failed to query job counts", http.StatusInternalServerError)
		log.Println("Failed to query job counts:", err)
		return
	}
	defer rows.Close()

	continentCounts := make(map[string]int)

	for rows.Next() {
		var continent string
		var count int
		err := rows.Scan(&continent, &count)
		if err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			log.Println("Failed to scan row:", err)
			return
		}
		continentCounts[continent] = count
	}

	helpers.RespondWithJSON(w, http.StatusOK, continentCounts)
}

func GetJobsByContinent(db_instance *sql.DB, w http.ResponseWriter, r *http.Request) {
	continent := chi.URLParam(r, "continent")
	if continent == "" {
		http.Error(w, "Continent parameter is required", http.StatusBadRequest)
		return
	}

	continent = strings.ToLower(continent)

	stmt, err := db_instance.Prepare("SELECT * FROM jobs WHERE LOWER(continent) = $1")
	if err != nil {
		http.Error(w, "Failed to prepare SQL statement", http.StatusInternalServerError)
		log.Println("Failed to prepare SQL statement:", err)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(continent)
	if err != nil {
		http.Error(w, "Failed to query jobs", http.StatusInternalServerError)
		log.Println("Failed to query jobs:", err)
		return
	}
	defer rows.Close()

	var jobs []types.Job

	for rows.Next() {
		var job types.Job
		err = rows.Scan(&job.JobID, &job.VideoID, &job.VideoUrl, &job.Title, &job.Description, &job.Country, &job.Category, &job.City, &job.Continent, &job.Company, &job.CompanyLogoUrl, &job.StartSalary, &job.EndSalary, &job.Currency, &job.PostedDate, &job.Deadline)
		if err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			log.Println("Failed to scan row:", err)
			return
		}

		jobs = append(jobs, job)
	}

	helpers.RespondWithJSON(w, http.StatusOK, jobs)
}

func GetJobsByCountry(db_instance *sql.DB, w http.ResponseWriter, r *http.Request) {
	country := chi.URLParam(r, "country")

	if country == "" {
		http.Error(w, "Country parameter is required", http.StatusBadRequest)
		return
	}

	stmt, err := db_instance.Prepare("SELECT * FROM jobs WHERE LOWER(country) = LOWER($1)")
	if err != nil {
		http.Error(w, "Failed to prepare SQL statement", http.StatusInternalServerError)
		log.Println("Failed to prepare SQL statement:", err)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(country)
	if err != nil {
		http.Error(w, "Failed to query jobs", http.StatusInternalServerError)
		log.Println("Failed to query jobs:", err)
		return
	}
	defer rows.Close()

	var jobs []types.Job

	for rows.Next() {
		var job types.Job

		err = rows.Scan(&job.JobID, &job.VideoID, &job.VideoUrl, &job.Title, &job.Description, &job.Country, &job.Category, &job.City, &job.Continent, &job.Company, &job.CompanyLogoUrl, &job.StartSalary, &job.EndSalary, &job.Currency, &job.PostedDate, &job.Deadline)
		if err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			log.Println("Failed to scan row:", err)
			return
		}

		jobs = append(jobs, job)
	}

	helpers.RespondWithJSON(w, http.StatusOK, jobs)
}
