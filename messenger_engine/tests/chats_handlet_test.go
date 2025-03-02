package tests

import (
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"testing"

	chathandler "messenger_engine/controllers/websocket_controller/handlers/chat_handler"
	chatcontroller "messenger_engine/controllers/chat_controller"
	errorhandler "messenger_engine/controllers/websocket_controller/handlers/error_handler"

	"github.com/gorilla/websocket"
)

// DummyChatController is a dummy implementation of the ChatControllerInterface.
// It provides a static JSON response for testing purposes.
type DummyChatController struct{}

// GetUserChats returns a dummy JSON response containing sample chat data.
func (d *DummyChatController) GetUserChats(userID int) ([]byte, error) {
	return []byte(`{"chats": ["chat1", "chat2"]}`), nil
}

// DummyErrorHandler is a dummy implementation of an error handler for WebSocket errors.
// In tests, it simply ignores errors.
type DummyErrorHandler struct{}

// HandleWebSocketError is a no-op implementation that satisfies the error handler interface.
func (d *DummyErrorHandler) HandleWebSocketError(err error, ws *websocket.Conn, format string, args ...interface{}) {
	// Intentionally left blank for testing purposes.
}

// TestChatsHandler_ServeHTTP tests the ServeHTTP method of the ChatsHandler.
// It verifies that the handler successfully upgrades the HTTP connection to a WebSocket,
// processes an incoming JSON message containing a user_id, and returns a response that
// includes the expected "chats" key.
func TestChatsHandler_ServeHTTP(t *testing.T) {
	// NOTE: The implementation below uses a concrete ChatController.
	// If your design allows dependency injection via an interface,
	// you can use &DummyChatController{} instead.
	dummyChatCtrl := &chatcontroller.ChatController{}
	handler := chathandler.NewChatsHandler(websocket.Upgrader{}, dummyChatCtrl)

	// Assign the error handler. You can also use a dummy if needed:
	// handler.ErrorHandler = &DummyErrorHandler{}
	handler.ErrorHandler = &errorhandler.ErrorHandler{}

	// Create a test HTTP server using the ChatsHandler.
	server := httptest.NewServer(handler)
	defer server.Close()

	// Convert the server URL from HTTP to WebSocket (ws) scheme.
	wsURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to parse server URL: %v", err)
	}
	wsURL.Scheme = "ws"

	// Establish a WebSocket connection to the test server.
	ws, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		t.Fatalf("failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	// Send a JSON request containing a user_id.
	request := map[string]interface{}{
		"user_id": 1,
	}
	if err := ws.WriteJSON(request); err != nil {
		t.Fatalf("failed to write JSON message: %v", err)
	}

	// Read the response message from the server.
	_, msg, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read message: %v", err)
	}

	// Unmarshal the JSON response into a map.
	var response map[string]interface{}
	if err := json.Unmarshal(msg, &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// Verify that the response contains the "chats" key.
	if _, ok := response["chats"]; !ok {
		t.Errorf("expected response to contain key 'chats', got: %v", response)
	}
}