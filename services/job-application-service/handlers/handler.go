package handlers

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/srishti13722/ai-hirehub/job-application-service/config"
	"github.com/srishti13722/ai-hirehub/job-application-service/kafka"
	"github.com/srishti13722/ai-hirehub/job-application-service/models"
)

func ApplyJob(c *fiber.Ctx) error {
	var application models.JobApplication
	userID := c.Locals("user_id")

	if err := c.BodyParser(&application); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	//fetch recruiter id
	fetchRecruiterId := `SELECT recruiter_id from jobs WHERE job_id = $1`
	err := config.DB.QueryRow(context.Background(), fetchRecruiterId, application.JobID).Scan(&application.RecruiterID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch recruiters details"})
	}

	//fetch job_seeker id
	fetchJobSeekerId := `SELECT jobseeker_id from job_seekers WHERE user_id = $1`
	err = config.DB.QueryRow(context.Background(), fetchJobSeekerId, userID).Scan(&application.JobSeekerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch jobseeker details"})
	}

	insertQuery := `INSERT INTO job_applications (job_id, jobseeker_id, recruiter_id) VALUES ($1, $2, $3)
	RETURNING application_id`

	//insert data into db

	err = config.DB.QueryRow(context.Background(), insertQuery, application.JobID,
		application.JobSeekerID, application.RecruiterID).Scan(&application.ApplicationID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create job application"})
	}

	event := models.ApplicationCreatedEvent{
		ApplicationID: application.ApplicationID,
		JobID:         application.JobID,
		RecruiterID:   application.RecruiterID,
		JobSeekerID:   application.JobSeekerID,
		AppliedAt:     time.Now(),
	}

	if err := kafka.Publish("job.application.created", event); err != nil {
		log.Println("Failed to publish Kafka event:", err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Job Application Created Successfully",
		"application_id": application.ApplicationID})
}

func GetMyApplications(c *fiber.Ctx) error {
	userID := c.Locals("user_id")

	// Fetch job_seeker_id
	var jobSeekerID string
	queryJobSeekerID := `SELECT jobseeker_id FROM job_seekers WHERE user_id = $1`
	err := config.DB.QueryRow(context.Background(), queryJobSeekerID, userID).Scan(&jobSeekerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch job seeker ID"})
	}

	query := `
	SELECT application_id, job_id, jobseeker_id, recruiter_id, status, applied_at, updated_at
	FROM job_applications WHERE jobseeker_id = $1
	`
	rows, err := config.DB.Query(context.Background(), query, jobSeekerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB Error: " + err.Error()})
	}
	defer rows.Close()

	var applications []models.JobApplication

	for rows.Next() {
		var app models.JobApplication
		err := rows.Scan(&app.ApplicationID, &app.JobID, &app.JobSeekerID, &app.RecruiterID, &app.Status, &app.AppliedAt, &app.UpdatedAt)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning application: " + err.Error()})
		}
		applications = append(applications, app)
	}

	return c.JSON(applications)
}

func GetJobApplications(c *fiber.Ctx) error {
	jobID := c.Params("job_id")
	userID := c.Locals("user_id")

	// Verify recruiter_id
	var recruiterID string
	fetchRecruiterID := `SELECT recruiter_id FROM recruiter WHERE user_id = $1`
	err := config.DB.QueryRow(context.Background(), fetchRecruiterID, userID).Scan(&recruiterID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch recruiter ID"})
	}

	query := `
	SELECT application_id, job_id, jobseeker_id, recruiter_id, status, applied_at, updated_at
	FROM job_applications
	WHERE job_id = $1 AND recruiter_id = $2
	`

	rows, err := config.DB.Query(context.Background(), query, jobID, recruiterID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB Error: " + err.Error()})
	}
	defer rows.Close()

	var applications []models.JobApplication

	for rows.Next() {
		var app models.JobApplication
		err := rows.Scan(&app.ApplicationID, &app.JobID, &app.JobSeekerID, &app.RecruiterID, &app.Status, &app.AppliedAt, &app.UpdatedAt)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning application: " + err.Error()})
		}
		applications = append(applications, app)
	}

	return c.JSON(applications)
}

func UpdateApplicationStatus(c *fiber.Ctx) error {
	applicationID := c.Params("id")
	userID := c.Locals("user_id")

	var statusUpdate struct {
		Status string `json:"status"`
	}

	if err := c.BodyParser(&statusUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Verify recruiter_id
	var recruiterID string
	fetchRecruiterID := `SELECT recruiter_id FROM recruiter WHERE user_id = $1`
	err := config.DB.QueryRow(context.Background(), fetchRecruiterID, userID).Scan(&recruiterID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch recruiter ID"})
	}

	// Update status
	query := `
	UPDATE job_applications SET status=$1, updated_at=NOW()
	WHERE application_id=$2 AND recruiter_id=$3
	`

	res, err := config.DB.Exec(context.Background(), query, statusUpdate.Status, applicationID, recruiterID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update status: " + err.Error()})
	}

	if res.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Application not found OR unauthorized"})
	}

	// Fetch job_id & jobseeker_id for this application
	var jobID, jobSeekerID string
	fetchQuery := `SELECT job_id, jobseeker_id FROM job_applications WHERE application_id = $1`
	err = config.DB.QueryRow(context.Background(), fetchQuery, applicationID).Scan(&jobID, &jobSeekerID)
	if err != nil {
		log.Println("Failed to fetch job or seeker ID for event:", err)
		return c.JSON(fiber.Map{"message": "Status updated, but failed to notify"})
	}

	// Prepare and publish Kafka event
	event := models.ApplicationStatusUpdatedEvent{
		ApplicationID: applicationID,
		JobID:         jobID,
		JobSeekerID:   jobSeekerID,
		NewStatus:     statusUpdate.Status,
		UpdatedAt:     time.Now(),
	}

	if err := kafka.Publish("job.application.status.updated", event); err != nil {
		log.Println("Failed to publish status update Kafka event:", err)
	}

	return c.JSON(fiber.Map{"message": "Application status updated"})
}

func DeleteApplication(c *fiber.Ctx) error {
	applicationID := c.Params("id")
	userID := c.Locals("user_id")

	// Fetch jobseeker_id
	var jobSeekerID string
	queryJobSeekerID := `SELECT jobseeker_id FROM job_seekers WHERE user_id = $1`
	err := config.DB.QueryRow(context.Background(), queryJobSeekerID, userID).Scan(&jobSeekerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch job seeker ID"})
	}

	// Delete query
	query := `DELETE FROM job_applications WHERE application_id=$1 AND jobseeker_id=$2`
	res, err := config.DB.Exec(context.Background(), query, applicationID, jobSeekerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete application: " + err.Error()})
	}

	if res.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Application not found OR unauthorized"})
	}

	return c.JSON(fiber.Map{"message": "Application deleted successfully"})
}
