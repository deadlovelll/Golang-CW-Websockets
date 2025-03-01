package errorhandler

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// ErrorHandler handles WebSocket errors.
type ErrorHandler struct{}

// NewErrorHandler creates and returns a new ErrorHandler instance.
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// HandleWebSocketError logs the error and sends an error message to the WebSocket client.
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
