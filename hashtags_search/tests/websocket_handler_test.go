package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"hashtags_search/handlers/websocket_handler"
)

// MockHashtagController is a mock implementation of the HashtagController interface for testing purposes.
// It simulates a database query by returning a fixed JSON response.
type MockHashtagController struct{}

// GetHashtags simulates a database query by returning a JSON-encoded slice of hashtag data.
// It returns a fixed response containing the provided hashtag and a dummy match count.
func (m *MockHashtagController) GetHashtags(hashtag string) ([]byte, error) {
	response := []map[string]interface{}{
		{"name": hashtag, "match_count": 10},
	}
	return json.Marshal(response)
}

// TestWebSocketHandler verifies that the WebSocketHandler correctly upgrades an HTTP connection to a WebSocket,
// processes incoming messages, and returns the expected response.
// It uses a test HTTP server to simulate a WebSocket connection, sends a test message,
// and validates that the response contains the correct hashtag data.
func TestWebSocketHandler(t *testing.T) {
	// Initialize the mock controller and WebSocket handler.
	dbCtrl := &MockHashtagController{}
	wsh := websockethandler.NewWebSocketHandler(dbCtrl)

	// Create a test HTTP server with the WebSocket handler.
	ts := httptest.NewServer(http.HandlerFunc(wsh.ServeHTTP))
	defer ts.Close()

	// Convert the HTTP test server URL to a WebSocket URL.
	wsURL := "ws" + ts.URL[4:]

	// Connect to the WebSocket server.
	wsConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("WebSocket connection error: %v", err)
	}
	defer wsConn.Close()

	// Create and send a test message with a sample query.
	testMessage := websockethandler.Message{Query: "example"}
	messageBytes, _ := json.Marshal(testMessage)
	if err := wsConn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
		t.Fatalf("Error sending WebSocket message: %v", err)
	}

	// Set a read deadline to avoid hanging if the response is not received.
	wsConn.SetReadDeadline(time.Now().Add(2 * time.Second))

	// Read the response from the WebSocket server.
	_, response, err := wsConn.ReadMessage()
	if err != nil {
		t.Fatalf("Error reading WebSocket response: %v", err)
	}

	// Parse the JSON response.
	var parsedResponse []map[string]interface{}
	if err := json.Unmarshal(response, &parsedResponse); err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}

	// Validate that the response contains the expected data.
	if len(parsedResponse) == 0 || parsedResponse[0]["name"] != "example" {
		t.Errorf("Unexpected response: %s", string(response))
	}
}
