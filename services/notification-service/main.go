package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"

	"github.com/srishti13722/ai-hirehub/notification-service/kafka"
	"github.com/srishti13722/ai-hirehub/notification-service/ws"
)

func main() {
	_ = godotenv.Load()

	app := fiber.New()

	log.Println("Kafka broker:", os.Getenv("KAFKA_BROKER"))

	// WebSocket route
	app.Get("/ws/:user_id", websocket.New(ws.WebSocketHandler))

	// Start Kafka consumers
	go kafka.ConsumeApplicationCreated()
	go kafka.ConsumeApplicationStatusUpdated()

	log.Println("Notification Service running on port 8086")
	app.Listen(":8086")
}
