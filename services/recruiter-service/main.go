package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/srishti13722/ai-hirehub/services/recruiter-service/config"
	"github.com/srishti13722/ai-hirehub/services/recruiter-service/handlers"
	"github.com/srishti13722/ai-hirehub/services/recruiter-service/middleware"
)

func main() {
	_ = godotenv.Load()
	config.ConnectDataBase()

	app := fiber.New()

	// Protected Routes
	recruiterGroup := app.Group("/recruiter", middleware.JWTMiddleware())
	recruiterGroup.Post("/create", handlers.CreateRecruiter)
	recruiterGroup.Get("/profile", handlers.GetRecruiter)
	recruiterGroup.Put("/update", handlers.UpdateRecruiter)
	recruiterGroup.Delete("/delete", handlers.DeleteRecruiter)

	log.Println("Recruiter Service running on port 8082")
	app.Listen(":8082")
}
