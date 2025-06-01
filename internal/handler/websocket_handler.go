package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)

// WebSocketHandler handles WebSocket connections
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ Upgrade failed:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true
	fmt.Println("✅ New WebSocket connection established")

	// Listen for messages from the WebSocket connection
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("❌ Error reading message:", err)
			break
		}
	}

	// Remove the client from the list when disconnected
	delete(clients, conn)
	fmt.Println("❌ WebSocket connection closed")
}

// WebSocketHandlerGin is a wrapper for WebSocketHandler to make it compatible with Gin
func WebSocketHandlerGin(c *gin.Context) {
	// Convert gin.Context to http.ResponseWriter and http.Request
	writer := c.Writer
	request := c.Request
	WebSocketHandler(writer, request)
}

// SendToClients sends a message to all connected WebSocket clients
func SendToClients(message string) {
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("❌ Error sending message:", err)
			client.Close()
			delete(clients, client)
		}
	}
}
