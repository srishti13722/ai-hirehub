package models

import "time"

type JobSeeker struct {
	JobSeekerID string    `json:"jobseeker_id"`
	UserID      string    `json:"user_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Experience  int       `json:"experience"`
	Skills      []string  `json:"skills"`
	ResumeURL   string    `json:"resume_url"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
