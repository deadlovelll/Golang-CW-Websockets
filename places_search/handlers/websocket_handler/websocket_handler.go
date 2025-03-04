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
// It contains a query string which can be a place name or a hashtag to search for places.
type Message struct {
	Query string `json:"query"` // The search query sent by the client.
}

// WebSocketHandler handles WebSocket connections and processes messages from clients.
type WebSocketHandler struct {
	upgrader websocket.Upgrader // WebSocket upgrader to upgrade HTTP connection to WebSocket.
	pCtrl    placecontroller.PlaceControllerInterface // A pointer to the PlaceController for querying place data.
}

// NewWebSocketHandler creates and returns a new WebSocketHandler instance.
// It initializes the WebSocket upgrader and sets up the place controller.
//
// Parameters:
//   - pCtrl: A pointer to the PlaceController that will handle place-related queries.
//
// Returns:
//   - A pointer to the newly created WebSocketHandler.
func NewWebSocketHandler(pCtrl placecontroller.PlaceControllerInterface) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			// CheckOrigin allows all connections. This can be customized based on security needs.
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		pCtrl: pCtrl,
	}
}

// ServeHTTP upgrades the HTTP connection to a WebSocket and processes incoming messages.
// This method listens for incoming WebSocket messages, processes them, and sends a response back to the client.
//
// Parameters:
//   - w: The HTTP response writer used to send the WebSocket response.
//   - r: The HTTP request that initiated the WebSocket connection.
func (wsh *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection.
	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to websocket: %v", err)
		return
	}
	defer conn.Close()

	// Continuously listen for incoming messages from the WebSocket connection.
	for {
		// Read the next WebSocket message.
		msgType, msgData, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading websocket message: %v", err)
			break
		}

		// Parse the incoming JSON message.
		var msg Message
		if err := json.Unmarshal(msgData, &msg); err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			continue
		}

		// Depending on the query, either search by place name or hashtag.
		var jsonResponse []byte
		if !strings.HasPrefix(msg.Query, "#") {
			// If the query doesn't start with a hashtag, search for place by name.
			jsonResponse, err = wsh.pCtrl.GetPlaceByName(msg.Query)
		} else {
			// If the query starts with a hashtag, search for places with that hashtag in the description.
			jsonResponse, err = wsh.pCtrl.GetPlaceWithHashtag(msg.Query)
		}

		// Handle any error that occurs while fetching data for the query.
		if err != nil {
			log.Printf("Error fetching data for query %s: %v", msg.Query, err)
			continue
		}

		// Send the JSON response back to the client.
		if err := conn.WriteMessage(msgType, jsonResponse); err != nil {
			log.Printf("Error sending websocket message: %v", err)
			break
		}

		log.Printf("Processed message: %s", msgData)
	}
}
