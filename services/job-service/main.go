package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/srishti13722/ai-hirehub/job-service/config"
	"github.com/srishti13722/ai-hirehub/job-service/handlers"
	"github.com/srishti13722/ai-hirehub/job-service/middleware"
)

func main() {
	_ = godotenv.Load()
	config.ConnectDataBase()

	app := fiber.New()

	app.Get("/jobs/list", handlers.GetJobList)
	app.Get("/job/:id", handlers.GetJobDetails)

	// Protected Routes
	jobGroup := app.Group("/job", middleware.JWTMiddleware())
	jobGroup.Post("/create", handlers.CreateJob)
	jobGroup.Put("/update/:id", handlers.UpdateJobDetails)
	jobGroup.Delete("/delete/:id", handlers.DeleteJob)
	


	log.Println("Job Service running on port 8084")
	app.Listen(":8084")
}
