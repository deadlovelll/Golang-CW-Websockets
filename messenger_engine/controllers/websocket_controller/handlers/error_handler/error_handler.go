package errorhandler

import (
	"fmt"
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

type WebSocketWriter interface {
    WriteJSON(v interface{}) error
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
// HandleWebSocketError writes the error message to the WebSocket connection.
func (e *ErrorHandler) HandleWebSocketError(err error, ws WebSocketWriter, format string, args ...interface{}) {
    message := map[string]string{
        "error": fmt.Sprintf(format, args...),
    }
    ws.WriteJSON(message)
}
