package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/srishti13722/ai-hirehub/services/recruiter-service/config"
	"github.com/srishti13722/ai-hirehub/services/recruiter-service/models"
)

func CreateRecruiter(c *fiber.Ctx) error {
	var recruiter models.Recruiter

	if err := c.BodyParser(&recruiter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Request: " + err.Error()})
	}

	recruiter.UserID = c.Locals("user_id").(string)

	query := `INSERT INTO recruiter (user_id, firstname, lastname, email, phone,
             organisation_name, designation,industry) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
			 RETURNING recruiter_id`

	err := config.DB.QueryRow(context.Background(), query, recruiter.UserID, recruiter.FirstName,
		recruiter.LastName, recruiter.Email, recruiter.Phone, recruiter.OrganisationName,
		recruiter.Designation, recruiter.Industry).Scan(&recruiter.RecruiterID)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB Error: Couldn't create recruiter: " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Recruiter Profile Created!", "recruiter_id": recruiter.RecruiterID})
}

func GetRecruiter(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var recruiter models.Recruiter

	recruiter.UserID = userID

	query := `SELECT recruiter_id, firstname, lastname, email, phone, organisation_name, 
	designation, industry, status, created_at, updated_at FROM recruiter WHERE user_id = $1`

	err := config.DB.QueryRow(context.Background(), query, userID).Scan(&recruiter.RecruiterID, &recruiter.FirstName,
		&recruiter.LastName, &recruiter.Email, &recruiter.Phone, &recruiter.OrganisationName, &recruiter.Designation, &recruiter.Industry,
		&recruiter.Status, &recruiter.CreatedAt, &recruiter.UpdatedAt)

	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB Error: Recruiter profile not found" + err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(recruiter)
}

func UpdateRecruiter(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var recruiter models.Recruiter

	if err := c.BodyParser(&recruiter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Request: " + err.Error()})
	}

	query := `
		UPDATE recruiter SET firstname=$1, lastname=$2, email=$3, phone=$4, organisation_name=$5, designation=$6, industry=$7, updated_at=$8
		WHERE user_id=$9
	`
	res, err := config.DB.Exec(context.Background(), query,
		recruiter.FirstName, recruiter.LastName, recruiter.Email, recruiter.Phone, recruiter.OrganisationName,
		recruiter.Designation, recruiter.Industry, time.Now(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating recruiter profile" + err.Error()})
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No recruiter profile found to update"})
	}

	return c.JSON(fiber.Map{"message": "Recruiter profile updated"})
}

func DeleteRecruiter(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	query := `DELETE FROM recruiter WHERE user_id=$1`
	_, err := config.DB.Exec(context.Background(), query, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting recruiter profile"})
	}
	return c.JSON(fiber.Map{"message": "Recruiter profile deleted"})
}
