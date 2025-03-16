package models

import "time"

type Recruiter struct {
	RecruiterID       string `json:"recruiter_id"`
	UserID            string `json:"user_id"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	OrganisationName  string `json:"organisation_name"`
	Designation       string `json:"designation"`
	Industry          string `json:"industry"`
	Status            string `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
