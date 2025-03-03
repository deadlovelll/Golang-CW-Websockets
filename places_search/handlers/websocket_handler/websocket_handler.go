package websockethandler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	placecontroller "places_search/controllers/place_controller"
)

// Message represents the structure of incoming WebSocket JSON messages.
//
// Fields:
//   - Query: A string representing the search query, which can be a place name or a hashtag.
type Message struct {
	Query string `json:"query"`
}

// WebSocketHandler manages WebSocket connections and processes messages.
//
// Fields:
//   - upgrader: Handles upgrading HTTP connections to WebSocket.
//   - pCtrl: A reference to PlaceController for handling place-related queries.
type WebSocketHandler struct {
	upgrader websocket.Upgrader
	pCtrl    *placecontroller.PlaceController
}

// NewWebSocketHandler initializes and returns a new instance of WebSocketHandler.
//
// Parameters:
//   - pCtrl: A pointer to an instance of PlaceController, used for processing queries.
//
// Returns:
//   - *WebSocketHandler: A pointer to the newly created WebSocketHandler.
func NewWebSocketHandler(pCtrl *placecontroller.PlaceController) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			// Allow all origins for WebSocket connections.
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		pCtrl: pCtrl,
	}
}

// ServeHTTP handles WebSocket requests by upgrading the HTTP connection and processing incoming messages.
//
// This function continuously listens for messages from the client, processes them based on the query type,
// and sends back the appropriate JSON response.
//
// Parameters:
//   - w: http.ResponseWriter used for responding to the client.
//   - r: *http.Request representing the incoming WebSocket request.
func (wsh *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to websocket: %v", err)
		return
	}
	defer conn.Close()

	for {
		// Read incoming WebSocket message.
		msgType, msgData, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading websocket message: %v", err)
			break
		}

		// Parse the JSON message into the Message struct.
		var msg Message
		if err := json.Unmarshal(msgData, &msg); err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			continue
		}

		// Determine query type: normal place name or hashtag search.
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

		// Send the response back to the WebSocket client.
		if err := conn.WriteMessage(msgType, jsonResponse); err != nil {
			log.Printf("Error sending websocket message: %v", err)
			break
		}

		log.Printf("Processed message: %s", msgData)
	}
}
