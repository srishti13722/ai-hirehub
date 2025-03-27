package ws

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

var (
	UserSockets = make(map[string]*websocket.Conn)
	mu          sync.RWMutex
)

func RegisterUser(userID string, conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	UserSockets[userID] = conn
}

func UnregisterUser(userID string) {
	mu.Lock()
	defer mu.Unlock()
	delete(UserSockets, userID)
}

func SendToUser(userID string, message []byte) error {
	mu.RLock()
	defer mu.RUnlock()

	if conn, ok := UserSockets[userID]; ok {
		return conn.WriteMessage(websocket.TextMessage, message)
	}
	return nil
}
