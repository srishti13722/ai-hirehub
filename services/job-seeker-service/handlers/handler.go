package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/srishti13722/ai-hirehub/services/job-seeker-service/config"
	"github.com/srishti13722/ai-hirehub/services/job-seeker-service/models"
)

func CreateJobSeeker(c *fiber.Ctx) error {
	var jobseeker models.JobSeeker

	if err := c.BodyParser(&jobseeker); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Request: " + err.Error()})
	}

	jobseeker.UserID = c.Locals("user_id").(string)

	query := `INSERT INTO job_seekers (user_id, firstname, lastname, email, phone, experience, skills, resume_url)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING jobseeker_id`

	err := config.DB.QueryRow(context.Background(), query,
		jobseeker.UserID,
		jobseeker.FirstName,
		jobseeker.LastName,
		jobseeker.Email,
		jobseeker.Phone,
		jobseeker.Experience,
		jobseeker.Skills,
		jobseeker.ResumeURL,
	).Scan(&jobseeker.JobSeekerID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB Error: Couldn't create profile: " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Job Seeker profile created!", "jobseeker_id": jobseeker.JobSeekerID})
}

func GetJobSeeker(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var jobseeker models.JobSeeker

	query := `SELECT jobseeker_id, user_id, firstname, lastname, email, phone, experience, skills, resume_url, status, created_at, updated_at 
			  FROM job_seekers WHERE user_id=$1`

	err := config.DB.QueryRow(context.Background(), query, userID).Scan(
		&jobseeker.JobSeekerID,
		&jobseeker.UserID,
		&jobseeker.FirstName,
		&jobseeker.LastName,
		&jobseeker.Email,
		&jobseeker.Phone,
		&jobseeker.Experience,
		&jobseeker.Skills,
		&jobseeker.ResumeURL,
		&jobseeker.Status,
		&jobseeker.CreatedAt,
		&jobseeker.UpdatedAt,
	)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Job Seeker profile not found: " + err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(jobseeker)
}

func UpdateJobSeeker(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var jobseeker models.JobSeeker

	if err := c.BodyParser(&jobseeker); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Request: " + err.Error()})
	}

	query := `
		UPDATE job_seekers SET firstname=$1, lastname=$2, email=$3, phone=$4, experience=$5, skills=$6, resume_url=$7, updated_at=$8
		WHERE user_id=$9
	`

	res, err := config.DB.Exec(context.Background(), query,
		jobseeker.FirstName,
		jobseeker.LastName,
		jobseeker.Email,
		jobseeker.Phone,
		jobseeker.Experience,
		jobseeker.Skills,
		jobseeker.ResumeURL,
		time.Now(),
		userID,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating profile: " + err.Error()})
	}

	if res.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No profile found to update"})
	}

	return c.JSON(fiber.Map{"message": "Job Seeker profile updated"})
}

func DeleteJobSeeker(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	query := `DELETE FROM job_seekers WHERE user_id=$1`
	res, err := config.DB.Exec(context.Background(), query, userID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting profile: " + err.Error()})
	}

	if res.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No profile found to delete"})
	}

	return c.JSON(fiber.Map{"message": "Job Seeker profile deleted"})
}

func UploadResume(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println("Error reading multipart form:", err)
	}
	fmt.Println("Multipart form received:", form)

	file, err := c.FormFile("resume")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No file uploaded"})
	}

	// Generate unique file name (to avoid conflicts)
	filename := fmt.Sprintf("%s_%s", userID, file.Filename)

	// Save file locally
	filePath := fmt.Sprintf("./uploads/%s", filename)
	err = c.SaveFile(file, filePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file: " + err.Error()})
	}

	// Update resume_url in DB
	query := `UPDATE job_seekers SET resume_url=$1, updated_at=$2 WHERE user_id=$3`
	res, err := config.DB.Exec(context.Background(), query, filePath, time.Now(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating resume URL: " + err.Error()})
	}

	if res.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No profile found to update"})
	}

	return c.JSON(fiber.Map{"message": "Resume uploaded successfully", "resume_url": filePath})
}
