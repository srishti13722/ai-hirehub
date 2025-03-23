package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/srishti13722/ai-hirehub/job-application-service/config"
	"github.com/srishti13722/ai-hirehub/job-application-service/handlers"
	"github.com/srishti13722/ai-hirehub/job-application-service/middleware"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using defaults")
	}

	config.ConnectDataBase()

	app := fiber.New()

	// For Job Seeker Routes
	jobSeekerGroup := app.Group("/seeker/application", middleware.JWTMiddleware("job_seeker", "admin"))
	jobSeekerGroup.Post("/apply", handlers.ApplyJob)
	jobSeekerGroup.Get("/my-applications", handlers.GetMyApplications)
	jobSeekerGroup.Delete("/delete/:id", handlers.DeleteApplication)

	// For Recruiter/Admin Routes
	recruiterGroup := app.Group("/recruiter/application", middleware.JWTMiddleware("recruiter", "admin"))
	recruiterGroup.Get("/job/:job_id", handlers.GetJobApplications)
	recruiterGroup.Put("/update-status/:id", handlers.UpdateApplicationStatus)

	log.Println("Job Application Service running on port 8085")
	app.Listen(":8085")
}
