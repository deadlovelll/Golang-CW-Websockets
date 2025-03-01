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
	upgrader websocket.Upgrader
	chatCtrl *chatcontroller.ChatController
	ErrorHandler  *ErrorHandler.ErrorHandler
}

// NewChatsHandler returns a new instance of ChatsHandler.
func NewChatsHandler(upgrader websocket.Upgrader, ctrl *chatcontroller.ChatController) *ChatsHandler {
	return &ChatsHandler{
		upgrader: upgrader,
		chatCtrl: ctrl,
	}
}

// ChatsMessage represents the expected JSON message.
type ChatsMessage struct {
	UserID int `json:"user_id"`
}

// ServeHTTP upgrades the connection and processes incoming messages.
func (h *ChatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := h.upgrader.Upgrade(w, r, nil)
	h.ErrorHandler.HandleWebSocketError(err, nil, "error %s when upgrading connection to websocket")

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			h.ErrorHandler.HandleWebSocketError(err, conn, "Error reading message %s")
			break
		}

		var msg ChatsMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			h.ErrorHandler.HandleWebSocketError(err, conn, "Error parsing message JSON: %s")
			continue
		}

		// Use the chat controller to fetch chat data.
		jsonData, err := h.chatCtrl.GetUserChats(msg.UserID)
		if err != nil {
			h.ErrorHandler.HandleWebSocketError(err, conn, "Error fetching users by ID: %s")
			return
		}

		if err := conn.WriteMessage(messageType, jsonData); err != nil {
			h.ErrorHandler.HandleWebSocketError(err, conn, "Error sending message: %v")
			break
		}

		log.Printf("Received message: %s", message)
	}
}
