// file: errorhandler_test.go
package tests

import (
	"errors"
	"testing"

	"messenger_engine/controllers/websocket_controller/handlers/error_handler"
)

// DummyWS is a mock WebSocket connection that records messages sent via WriteJSON.
type DummyWS struct {
	messages []interface{}
}

// WriteJSON simulates writing a JSON message to a WebSocket connection
// by appending it to the messages slice.
func (d *DummyWS) WriteJSON(v interface{}) error {
	d.messages = append(d.messages, v)
	return nil
}

// TestHandleWebSocketError_WithWS verifies that HandleWebSocketError
// correctly sends an error message to a WebSocket connection.
func TestHandleWebSocketError_WithWS(t *testing.T) {
	eh := errorhandler.NewErrorHandler()

	// Use the custom DummyWS, not *websocket.Conn
	dummyWS := &DummyWS{}
	testErr := errors.New("test error")

	// Call the error handler, passing the dummyWS (of type *DummyWS)
	eh.HandleWebSocketError(testErr, dummyWS, "Error: %s", "failed")

	// Check that a message was recorded in the dummyWS.messages slice
	if len(dummyWS.messages) == 0 {
		t.Error("expected error message to be written to websocket, got none")
	} else {
		// Validate that the message is a map with an "error" field
		msg, ok := dummyWS.messages[0].(map[string]string)
		if !ok {
			t.Error("expected message to be a map[string]string")
		}
		if msg["error"] == "" {
			t.Error("expected non-empty error field in the message")
		}
	}
}

// TestHandleWebSocketError_NilWS ensures that HandleWebSocketError does not panic
// when given a nil WebSocket connection.
func TestHandleWebSocketError_NilWS(t *testing.T) {
	eh := errorhandler.NewErrorHandler()

	// Passing a nil WebSocket should not cause a panic.
	eh.HandleWebSocketError(errors.New("test error"), nil, "Error: %s", "failed")
}
