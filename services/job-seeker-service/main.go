package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/srishti13722/ai-hirehub/job-seeker-service/config"
	"github.com/srishti13722/ai-hirehub/job-seeker-service/handlers"
	"github.com/srishti13722/ai-hirehub/job-seeker-service/middleware"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using defaults")
	}

	// Connect DB
	config.ConnectDataBase()
	config.RunMigrations()

	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024,
		StreamRequestBody: true, 
	})

	// Protected Routes Group (Only job_seekers & admins allowed)
	jobseekerGroup := app.Group("/jobseeker", middleware.JWTMiddleware())

	jobseekerGroup.Post("/create", handlers.CreateJobSeeker)
	jobseekerGroup.Get("/profile", handlers.GetJobSeeker)
	jobseekerGroup.Put("/update", handlers.UpdateJobSeeker)
	jobseekerGroup.Delete("/delete", handlers.DeleteJobSeeker)
	jobseekerGroup.Post("/upload-resume", handlers.UploadResume)

	// Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "Job Seeker Service running"})
	})

	log.Println("Job Seeker Service running on port 8083")
	app.Listen(":8083")
}
