package websockethandler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"hashtags_search/controllers/hashtag_controller"
)

// Message represents the structure of incoming WebSocket messages.
type Message struct {
	Query string `json:"query"`
}

// WebSocketHandler manages WebSocket connections and message handling.
type WebSocketHandler struct {
	upgrader websocket.Upgrader
	dbCtrl   hashtagcontroller.HashtagProvider
}

// NewWebSocketHandler initializes a new WebSocket handler with the given database controller.
func NewWebSocketHandler(dbCtrl hashtagcontroller.HashtagProvider) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		dbCtrl:   dbCtrl,
	}
}


// ServeHTTP upgrades an HTTP connection to a WebSocket and processes messages.
func (wsh *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	for {
		msgType, msgData, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading WebSocket message: %v", err)
			break
		}

		var msg Message
		if err := json.Unmarshal(msgData, &msg); err != nil {
			log.Printf("Error parsing JSON message: %v", err)
			continue
		}

		jsonData, err := wsh.dbCtrl.GetHashtags(msg.Query)
		if err != nil {
			log.Printf("Error fetching hashtags for query '%s': %v", msg.Query, err)
			return
		}

		log.Printf("Received query: %s; Responding with: %s", msg.Query, string(jsonData))

		if err := conn.WriteMessage(msgType, jsonData); err != nil {
			log.Printf("Error sending WebSocket message: %v", err)
			break
		}
	}
}
