package models

import "time"

type ApplicationCreatedEvent struct {
	ApplicationID string    `json:"application_id"`
	JobID         string    `json:"job_id"`
	RecruiterID   string    `json:"recruiter_id"`
	JobSeekerID   string    `json:"jobseeker_id"`
	AppliedAt     time.Time `json:"applied_at"`
}

type ApplicationStatusUpdatedEvent struct {
	ApplicationID string    `json:"application_id"`
	JobID         string    `json:"job_id"`
	JobSeekerID   string    `json:"jobseeker_id"`
	NewStatus     string    `json:"new_status"`
	UpdatedAt     time.Time `json:"updated_at"`
}