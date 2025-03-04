package tests

import (
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"places_search/handlers/websocket_handler"
)

// DummyPlaceController is a fake implementation of PlaceControllerInterface for testing.
// It provides dummy responses for both GetPlaceByName and GetPlaceWithHashtag methods.
type DummyPlaceController struct{}

// GetPlaceByName returns a dummy JSON response indicating a search by name.
// This simulates the behavior of a real PlaceController when a name query is executed.
//
// Parameters:
//   - query: The search query string.
//
// Returns:
//   - []byte: The JSON encoded response.
//   - error: An error if JSON marshaling fails.
func (d *DummyPlaceController) GetPlaceByName(query string) ([]byte, error) {
	response := map[string]string{
		"result": "searched by name: " + query,
	}
	return json.Marshal(response)
}

// GetPlaceWithHashtag returns a dummy JSON response indicating a search by hashtag.
// This simulates the behavior of a real PlaceController when a hashtag query is executed.
//
// Parameters:
//   - query: The hashtag query string.
//
// Returns:
//   - []byte: The JSON encoded response.
//   - error: An error if JSON marshaling fails.
func (d *DummyPlaceController) GetPlaceWithHashtag(query string) ([]byte, error) {
	response := map[string]string{
		"result": "searched by hashtag: " + query,
	}
	return json.Marshal(response)
}

// TestWebSocketHandler_GetPlaceByName tests that a plain query (without a hashtag)
// returns the expected JSON response via the WebSocketHandler.
func TestWebSocketHandler_GetPlaceByName(t *testing.T) {
	// Create a dummy place controller.
	fakeCtrl := &DummyPlaceController{}

	// Create the WebSocketHandler using the dummy controller.
	wsHandler := websockethandler.NewWebSocketHandler(fakeCtrl)

	// Start a test HTTP server using the WebSocketHandler.
	server := httptest.NewServer(wsHandler)
	defer server.Close()

	// Convert the test server's URL to a WebSocket URL.
	u, err := url.Parse(server.URL)
	assert.NoError(t, err)
	wsURL := "ws://" + u.Host

	// Dial the WebSocket server.
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer conn.Close()

	// Prepare and send a test message (without a hashtag).
	msg := websockethandler.Message{Query: "test"}
	msgBytes, err := json.Marshal(msg)
	assert.NoError(t, err)
	err = conn.WriteMessage(websocket.TextMessage, msgBytes)
	assert.NoError(t, err)

	// Read the response from the server.
	_, respBytes, err := conn.ReadMessage()
	assert.NoError(t, err)

	// Verify the response content.
	var respData map[string]string
	err = json.Unmarshal(respBytes, &respData)
	assert.NoError(t, err)
	assert.Equal(t, "searched by name: test", respData["result"])
}

// TestWebSocketHandler_GetPlaceWithHashtag tests that a query starting with a hashtag
// returns the expected JSON response via the WebSocketHandler.
func TestWebSocketHandler_GetPlaceWithHashtag(t *testing.T) {
	// Create a dummy place controller.
	fakeCtrl := &DummyPlaceController{}

	// Create the WebSocketHandler using the dummy controller.
	wsHandler := websockethandler.NewWebSocketHandler(fakeCtrl)

	// Start a test HTTP server using the WebSocketHandler.
	server := httptest.NewServer(wsHandler)
	defer server.Close()

	// Convert the test server's URL to a WebSocket URL.
	u, err := url.Parse(server.URL)
	assert.NoError(t, err)
	wsURL := "ws://" + u.Host

	// Dial the WebSocket server.
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer conn.Close()

	// Prepare and send a test message (with a hashtag).
	msg := websockethandler.Message{Query: "#test"}
	msgBytes, err := json.Marshal(msg)
	assert.NoError(t, err)
	err = conn.WriteMessage(websocket.TextMessage, msgBytes)
	assert.NoError(t, err)

	// Read the response from the server.
	_, respBytes, err := conn.ReadMessage()
	assert.NoError(t, err)

	// Verify the response content.
	var respData map[string]string
	err = json.Unmarshal(respBytes, &respData)
	assert.NoError(t, err)
	assert.Equal(t, "searched by hashtag: #test", respData["result"])
}
