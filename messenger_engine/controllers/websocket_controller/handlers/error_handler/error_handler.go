package errorhandler

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// ErrorHandler provides methods for handling errors.
type ErrorHandler struct{}

// HandleWebSocketError logs the error and sends an error message to the WebSocket client.
func (eh *ErrorHandler) HandleWebSocketError(err error, ws *websocket.Conn, message string) {
	if err != nil {
		fmt.Printf(message+"\n", err)
		if ws != nil {
			ws.WriteJSON(map[string]string{"error": fmt.Sprintf(message, err)})
		}
	}
}
