package websockethandler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	placecontroller "places_search/controllers/place_controller"
)

// Message represents the incoming WebSocket JSON message.
type Message struct {
	Query string `json:"query"`
}

// WebSocketHandler handles WebSocket connections.
type WebSocketHandler struct {
	upgrader websocket.Upgrader
	pCtrl   *placecontroller.PlaceController
}

// NewWebSocketHandler creates and returns a new WebSocketHandler.
func NewWebSocketHandler(pCtrl *placecontroller.PlaceController) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		pCtrl: pCtrl,
	}
}

// ServeHTTP upgrades the HTTP connection to a WebSocket and processes incoming messages.
func (wsh *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to websocket: %v", err)
		return
	}
	defer conn.Close()

	for {
		msgType, msgData, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading websocket message: %v", err)
			break
		}

		// Parse the JSON message.
		var msg Message
		if err := json.Unmarshal(msgData, &msg); err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			continue
		}

		// Check if the query starts with "#" to decide which controller method to use.
		var jsonResponse []byte
		if !strings.HasPrefix(msg.Query, "#") {
			jsonResponse, err = wsh.pCtrl.GetPlaceByName(msg.Query)
		} else {
			jsonResponse, err = wsh.pCtrl.GetPlaceWithHashtag(msg.Query)
		}
		if err != nil {
			log.Printf("Error fetching data for query %s: %v", msg.Query, err)
			continue
		}

		// Write the JSON response back to the client.
		if err := conn.WriteMessage(msgType, jsonResponse); err != nil {
			log.Printf("Error sending websocket message: %v", err)
			break
		}

		log.Printf("Processed message: %s", msgData)
	}
}
