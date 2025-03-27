package ws

import (
	"log"

	"github.com/gofiber/websocket/v2"
)

func WebSocketHandler(c *websocket.Conn) {
	userID := c.Params("user_id")

	RegisterUser(userID, c)
	log.Printf("User %s connected", userID)

	defer func() {
		UnregisterUser(userID)
		log.Printf("User %s disconnected", userID)
		c.Close()
	}()

	for {
		if _, _, err := c.ReadMessage(); err != nil {
			break
		}
	}
}
