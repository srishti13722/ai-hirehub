package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/srishti13722/ai-hirehub/auth-service/config"
	"github.com/srishti13722/ai-hirehub/auth-service/handlers"
	"github.com/srishti13722/ai-hirehub/auth-service/middleware"
)

func main(){
	_ = godotenv.Load()
	config.ConnectDataBase()

	app := fiber.New()

	// Auth routes
	app.Post("/signup", handlers.SignUp)
	app.Post("/login", handlers.Login)

	// Protected Routes
	adminRoutes := app.Group("/admin", middleware.RoleBasedMiddleware("admin"))
	adminRoutes.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome, Admin!"})
	})

	recruiterRoutes := app.Group("/recruiter", middleware.RoleBasedMiddleware("recruiter"))
	recruiterRoutes.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome, Recruiter!"})
	})

	log.Println("Auth Service running on port 8081")
	app.Listen(":8081")
}