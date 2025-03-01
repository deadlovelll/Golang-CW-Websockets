package errorhandler

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// ErrorHandler provides methods for handling WebSocket errors.
type ErrorHandler struct{}

// NewErrorHandler creates and returns a new instance of ErrorHandler.
//
// Returns:
//   - A pointer to an initialized ErrorHandler.
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// HandleWebSocketError logs the given error and sends a formatted error message to the WebSocket client.
//
// Parameters:
//   - err: The error encountered (if nil, the function returns without any action).
//   - ws: A pointer to the WebSocket connection to which the error message will be sent.
//   - format: A format string for the error message (similar to fmt.Sprintf).
//   - args: Additional arguments to be formatted into the error message.
//
// Behavior:
//   - If `err` is nil, the function does nothing.
//   - Logs the error message with the provided format.
//   - Attempts to send a JSON-encoded error response to the WebSocket client.
func (eh *ErrorHandler) HandleWebSocketError(err error, ws *websocket.Conn, format string, args ...interface{}) {
	if err == nil {
		return
	}

	// Format the error message
	message := fmt.Sprintf(format, args...)
	log.Printf("WebSocket Error: %s - %v", message, err)

	// Send error response to the WebSocket client, if available
	if ws != nil {
		_ = ws.WriteJSON(map[string]string{"error": message})
	}
}
