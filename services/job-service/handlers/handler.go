package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/srishti13722/ai-hirehub/job-service/config"
	"github.com/srishti13722/ai-hirehub/job-service/models"
)

func CreateJob(c *fiber.Ctx) error {
	var job models.Job

	if err := c.BodyParser(&job); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	userID := c.Locals("user_id")

	//fetch recruiter id
	fetchRecruiterId := `SELECT recruiter_id FROM recruiter WHERE user_id=$1`
	err := config.DB.QueryRow(context.Background(), fetchRecruiterId, userID).Scan(&job.RecruiterID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "couldn't fetch recruiter id"})
	}

	//insert job data into db

	insertJob := `INSERT INTO jobs(recruiter_id, organisation_name, job_title, job_description, job_location,
	salary, skills_required, vacancy) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING job_id`

	err = config.DB.QueryRow(context.Background(), insertJob, job.RecruiterID, job.OrganisationName,
		job.JobTitle, job.JobDescription, job.JobLocation, job.Salary, job.SkillsRequired, job.Vacancy).Scan(&job.JobID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failled to create the job"})
	}

	return c.JSON(fiber.Map{"message": "successfully created the job", "job_id": job.JobID})
}

func GetJobList(c *fiber.Ctx) error {
	location := c.Query("location")
	skills := c.Query("skills")
	minSalary := c.Query("min_salary")
	maxSalary := c.Query("max_salary")
	jobTitle := c.Query("job_title")
	limit := c.Query("limit", "10")
	offset := c.Query("offset", "0")

	// Build Redis key from all filters
	cacheKey := fmt.Sprintf("jobs:location=%s|skills=%s|minSalary=%s|maxSalary=%s|title=%s|limit=%s|offset=%s",
		location, skills, minSalary, maxSalary, jobTitle, limit, offset,
	)

	// Try to get from Redis cache
	if cached, err := config.RedisClient.Get(config.Ctx, cacheKey).Result(); err == nil {
		var cachedJobs []models.Job
		if err := json.Unmarshal([]byte(cached), &cachedJobs); err == nil {
			return c.JSON(fiber.Map{"cached": true, "jobs": cachedJobs})
		}
	}

	// Build SQL query
	query := `
	SELECT job_id, recruiter_id, organisation_name, job_title, job_description, job_location,
	salary, skills_required, vacancy, status, created_at, updated_at
	FROM jobs WHERE status = 'active'
	`

	params := []any{}
	paramCounter := 1

	if location != "" {
		query += fmt.Sprintf(" AND job_location ILIKE $%d", paramCounter)
		params = append(params, "%"+location+"%")
		paramCounter++
	}

	if skills != "" {
		query += fmt.Sprintf(" AND skills_required @> $%d", paramCounter)
		params = append(params, strings.Split(skills, ","))
		paramCounter++
	}

	if minSalary != "" {
		query += fmt.Sprintf(" AND salary >= $%d", paramCounter)
		params = append(params, minSalary)
		paramCounter++
	}

	if maxSalary != "" {
		query += fmt.Sprintf(" AND salary <= $%d", paramCounter)
		params = append(params, maxSalary)
		paramCounter++
	}

	if jobTitle != "" {
		query += fmt.Sprintf(" AND job_title ILIKE $%d", paramCounter)
		params = append(params, "%"+jobTitle+"%")
		paramCounter++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", paramCounter, paramCounter+1)
	params = append(params, limit, offset)

	// Execute query
	rows, err := config.DB.Query(context.Background(), query, params...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB Error: " + err.Error()})
	}
	defer rows.Close()

	var jobs []models.Job

	for rows.Next() {
		var job models.Job
		err := rows.Scan(
			&job.JobID,
			&job.RecruiterID,
			&job.OrganisationName,
			&job.JobTitle,
			&job.JobDescription,
			&job.JobLocation,
			&job.Salary,
			&job.SkillsRequired,
			&job.Vacancy,
			&job.Status,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning job: " + err.Error()})
		}
		jobs = append(jobs, job)
	}

	// Save to Redis cache
	if jobBytes, err := json.Marshal(jobs); err == nil {
		config.RedisClient.Set(config.Ctx, cacheKey, jobBytes, 60*time.Second)
	}

	return c.JSON(fiber.Map{"cached": false, "jobs": jobs})
}

func GetJobDetails(c *fiber.Ctx) error {
	jobID := c.Params("id")

	query := `SELECT job_id, recruiter_id, organisation_name, job_title, job_description,
	 job_location, salary, skills_required, vacancy, status, created_at, updated_at
		FROM jobs WHERE job_id = $1`

	var job models.Job

	err := config.DB.QueryRow(context.Background(), query, jobID).Scan(&job.JobID, &job.RecruiterID, &job.OrganisationName,
		&job.JobTitle, &job.JobDescription, &job.JobLocation, &job.Salary, &job.SkillsRequired,
		&job.Vacancy, &job.Status, &job.CreatedAt, &job.UpdatedAt)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch job details" + err.Error()})
	}

	return c.JSON(job)
}

func UpdateJobDetails(c *fiber.Ctx) error {
	jobID := c.Params("id")
	userID := c.Locals("user_id")

	var recruiterID string

	//fetch recruiter id
	fetchRecruiterId := `SELECT recruiter_id FROM recruiter WHERE user_id=$1`
	err := config.DB.QueryRow(context.Background(), fetchRecruiterId, userID).Scan(&recruiterID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "couldn't fetch recruiter id"})
	}

	var job models.Job

	if err := c.BodyParser(&job); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	query := `UPDATE jobs SET organisation_name=$1, job_title=$2, job_description=$3, job_location=$4,
	 salary=$5, skills_required=$6, vacancy=$7, updated_at=$8
		WHERE job_id=$9 AND recruiter_id=$10`

	res, err := config.DB.Exec(context.Background(), query, job.OrganisationName, job.JobTitle,
		job.JobDescription, job.JobLocation, job.Salary, job.SkillsRequired, job.Vacancy, time.Now(),
		jobID, recruiterID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update job: " + err.Error()})
	}

	if res.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No job found OR unauthorized"})
	}

	return c.JSON(fiber.Map{"message": "Job details updated successfully!"})
}

func DeleteJob(c *fiber.Ctx) error {
	jobID := c.Params("id")
	userID := c.Locals("user_id")

	var recruiterID string

	//fetch recruiter id
	fetchRecruiterId := `SELECT recruiter_id FROM recruiter WHERE user_id=$1`
	err := config.DB.QueryRow(context.Background(), fetchRecruiterId, userID).Scan(&recruiterID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "couldn't fetch recruiter id"})
	}

	query := `DELETE FROM jobs WHERE job_id = $1 and recruiter_id = $2`

	res, err := config.DB.Exec(context.Background(), query, jobID, recruiterID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete job: " + err.Error()})
	}

	if res.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No job found OR unauthorized"})
	}

	return c.JSON(fiber.Map{"message": "Job deleted successfully!"})
}
