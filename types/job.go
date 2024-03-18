package types

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
