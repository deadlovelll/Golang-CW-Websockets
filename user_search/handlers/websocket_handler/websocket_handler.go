package websockethandler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	usercontroller "user_search/controllers/user_controller"

	"github.com/gorilla/websocket"
)

// Message represents an incoming message from the client.
type Message struct {
	Query string `json:"query"`
}

// WebSocketHandler handles WebSocket connections.
type WebSocketHandler struct {
	upgrader websocket.Upgrader
	userCtrl usercontroller.UserControllerInterface
}

// NewWebSocketHandler creates a new instance of WebSocketHandler.
func NewWebSocketHandler(userCtrl usercontroller.UserControllerInterface) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			// Allow all origins for simplicity. Adjust as needed for production.
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		userCtrl: userCtrl,
	}
}

// ServeHTTP handles HTTP requests and upgrades the connection to a WebSocket.
func (wsh *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection.
	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	for {
		// Read message from WebSocket connection.
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Parse the JSON message.
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error parsing JSON: %v", err)
			continue
		}

		// If the query can be converted to an integer, assume it's a UserID.
		if userID, err := strconv.Atoi(msg.Query); err == nil {
			jsonData, err := wsh.userCtrl.GetUsers(userID)
			if err != nil {
				log.Printf("Error fetching users by ID: %v", err)
				continue
			}
			if err := conn.WriteMessage(messageType, jsonData); err != nil {
				log.Printf("Error sending message: %v", err)
				break
			}
			log.Printf("Message received: %s", message)
		} else {
			// Otherwise, assume the query is a username.
			username := strings.TrimSpace(msg.Query)
			if username == "" {
				log.Println("Invalid username")
				continue
			}
			jsonData, err := wsh.userCtrl.GetUsersByUsername(username)
			if err != nil {
				log.Printf("Error fetching users by username: %v", err)
				continue
			}
			if err := conn.WriteMessage(messageType, jsonData); err != nil {
				log.Printf("Error sending message: %v", err)
				break
			}
			log.Printf("Message received: %s", message)
		}
	}
}
