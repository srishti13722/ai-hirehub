package models

import "time"

type JobApplication struct {
	ApplicationID string    `json:"application_id"`
	JobID         string    `json:"job_id"`
	JobSeekerID   string    `json:"jobseeker_id"`
	RecruiterID   string    `json:"recruiter_id"`
	Status        string    `json:"status"`
	AppliedAt     time.Time `json:"applied_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
