package broadcastcontroller

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"

	"messenger_engine/models/message"
	"messenger_engine/controllers/message_controller"
)

// Broadcast manages WebSocket clients and handles message broadcasting.
type Broadcast struct {
	Clients         map[*websocket.Conn]bool // A map of connected WebSocket clients.
	mu              sync.Mutex               // Mutex to ensure concurrent safety.
	Broadcast       chan message.FinalMessage       // Channel for broadcasting messages.
	RepliesBroadcast chan message.FinalMessageReply // Channel for broadcasting reply messages.
}

// NewBroadcaster initializes and returns a new Broadcast instance.
//
// Returns:
//   - A pointer to a newly created Broadcast instance.
func NewBroadcaster() *Broadcast {
	return &Broadcast{
		Clients:         make(map[*websocket.Conn]bool),
		Broadcast:       make(chan message.FinalMessage),
		RepliesBroadcast: make(chan message.FinalMessageReply),
	}
}

// RegisterClient adds a new WebSocket client to the broadcaster.
//
// Parameters:
//   - client: The WebSocket connection to be registered.
func (b *Broadcast) RegisterClient(client *websocket.Conn) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Clients[client] = true
}

// RemoveClient removes a WebSocket client from the broadcaster and closes its connection.
//
// Parameters:
//   - client: The WebSocket connection to be removed.
func (b *Broadcast) RemoveClient(client *websocket.Conn) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, exists := b.Clients[client]; exists {
		client.Close()
		delete(b.Clients, client)
	}
}

// HandleMessages listens for incoming messages on the Broadcast channel
// and sends them to all registered WebSocket clients.
//
// Parameters:
//   - mmc: A pointer to the MessageController that handles messages.
func (b *Broadcast) HandleMessages(mmc *messagecontroller.MessageController) {
	for msg := range b.Broadcast {
		b.mu.Lock()
		for client := range b.Clients {
			if err := client.WriteJSON(msg); err != nil {
				log.Printf("Error sending message to client: %v", err)
				b.RemoveClient(client)
			}
		}
		b.mu.Unlock()
	}
}

// BroadcastMessage sends a message to all registered WebSocket clients.
//
// Parameters:
//   - msg: The FinalMessage struct containing the message to broadcast.
func (b *Broadcast) BroadcastMessage(msg message.FinalMessage) {
	select {
	case b.Broadcast <- msg:
	default:
		log.Println("Broadcast channel full, dropping message")
	}
}

// BroadcastReplyMessage sends a reply message to all registered WebSocket clients.
//
// Parameters:
//   - msg: The FinalMessageReply struct containing the reply message to broadcast.
func (b *Broadcast) BroadcastReplyMessage(msg message.FinalMessageReply) {
	select {
	case b.RepliesBroadcast <- msg:
	default:
		log.Println("Reply broadcast channel full, dropping message")
	}
}
