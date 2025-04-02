package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/srishti13722/ai-hirehub/ai-service/config"
	"github.com/srishti13722/ai-hirehub/ai-service/handler"
)

func main() {
	_ = godotenv.Load()

	config.InitOpenAI()

	app := fiber.New()

	app.Post("/parse-resume", handler.ParseResume)
	app.Post("/generate-cover-letter", handler.GenerateCoverLetter)
	app.Post("/rank-candidates", handler.RankCandidates)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "AI Service running"})
	})

	log.Println("🚀 AI Service running on port 8087")
	app.Listen(":8087")
}
