package models

import "time"

type Job struct {
	JobID            string    `json:"job_id"`
	RecruiterID      string    `json:"recruiter_id"`
	OrganisationName string    `json:"organisation_name"`
	JobTitle         string    `json:"job_title"`
	JobDescription   string    `json:"job_description"`
	JobLocation      string    `json:"job_location"`
	Salary           int       `json:"salary"`
	SkillsRequired   []string  `json:"skills_required"`
	Vacancy          int       `json:"vacancy"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
