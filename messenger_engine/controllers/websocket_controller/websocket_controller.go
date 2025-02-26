package websocketcontroller

import (
	"fmt"

	Messages "messenger_engine/models/message"

	"github.com/gorilla/websocket"
)

// ----------------------------------------------------------------------
// WebSocket Broadcaster
// ----------------------------------------------------------------------

// MessageBroadcaster manages WebSocket clients and message broadcasting.
type MessageBroadcaster struct {
	Clients          map[*websocket.Conn]bool
	Broadcast        chan Messages.FinalMessage
	RepliesBroadcast chan Messages.FinalMessageReply
}

// NewMessageBroadcaster initializes and returns a new MessageBroadcaster.
func NewMessageBroadcaster() *MessageBroadcaster {
	return &MessageBroadcaster{
		Clients:          make(map[*websocket.Conn]bool),
		Broadcast:        make(chan Messages.FinalMessage),
		RepliesBroadcast: make(chan Messages.FinalMessageReply),
	}
}

// HandleMessages continuously reads messages from the broadcast channel and sends them to connected clients.
func (mb *MessageBroadcaster) HandleMessages() {
	for msg := range mb.Broadcast {
		for client := range mb.Clients {
			if err := client.WriteJSON(msg); err != nil {
				fmt.Printf("Error sending message: %v\n", err)
				client.Close()
				delete(mb.Clients, client)
			}
		}
	}
}