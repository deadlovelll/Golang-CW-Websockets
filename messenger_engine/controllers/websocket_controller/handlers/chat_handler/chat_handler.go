package chathandler

import (
	"encoding/json"
	"log"
	"net/http"

	"messenger_engine/controllers/chat_controller"
	ErrorHandler "messenger_engine/controllers/websocket_controller/handlers/error_handler"

	"github.com/gorilla/websocket"
)

// ChatsHandler handles WebSocket connections for retrieving user chats.
type ChatsHandler struct {
	upgrader     websocket.Upgrader  // WebSocket upgrader to upgrade the HTTP connection
	chatCtrl     *chatcontroller.ChatController // Controller for managing chat-related operations
	ErrorHandler *ErrorHandler.ErrorHandler // Error handler for managing WebSocket-related errors
}

// NewChatsHandler initializes a new ChatsHandler instance with the given WebSocket upgrader and chat controller.
// Returns a pointer to a new ChatsHandler.
func NewChatsHandler(upgrader websocket.Upgrader, ctrl *chatcontroller.ChatController) *ChatsHandler {
	return &ChatsHandler{
		upgrader: upgrader,
		chatCtrl: ctrl,
	}
}

// ChatsMessage represents the expected JSON message structure containing the user ID.
type ChatsMessage struct {
	UserID int `json:"user_id"` // The user ID for retrieving their chats
}

// ServeHTTP upgrades the connection to WebSocket, processes incoming chat requests,
// and sends back the relevant chat data to the client.
// It listens for incoming messages, parses them, and retrieves chats for the given user.
func (h *ChatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Allow WebSocket connections from any origin
	h.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		// Handle WebSocket upgrade error
		h.ErrorHandler.HandleWebSocketError(err, nil, "Error upgrading connection to WebSocket: %s", err)
		return
	}
	defer conn.Close() // Ensure the connection is closed when the function exits

	// Process incoming messages in a loop
	for {
		// Read message from the WebSocket connection
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			// Handle error in reading the message
			h.ErrorHandler.HandleWebSocketError(err, conn, "Error reading message: %s", err)
			break
		}

		// Parse the received message into ChatsMessage
		var msg ChatsMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			// Handle error in parsing the JSON message
			h.ErrorHandler.HandleWebSocketError(err, conn, "Error parsing message JSON: %s", err)
			continue
		}

		// Use the chat controller to fetch the user's chats
		jsonData, err := h.chatCtrl.GetUserChats(msg.UserID)
		if err != nil {
			// Handle error in fetching user chats
			h.ErrorHandler.HandleWebSocketError(err, conn, "Error fetching users by ID: %s", err)
			return
		}

		// Send the fetched chat data back to the client
		if err := conn.WriteMessage(messageType, jsonData); err != nil {
			// Handle error in sending message to client
			h.ErrorHandler.HandleWebSocketError(err, conn, "Error sending message: %v", err)
			break
		}

		// Log the received message for debugging purposes
		log.Printf("Received message: %s", message)
	}
}