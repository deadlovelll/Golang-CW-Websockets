package broadcastcontroller

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"

	"messenger_engine/models/message"
	"messenger_engine/controllers/message_controller"
)

// Broadcaster manages WebSocket clients and message broadcasting.
type Broadcast struct {
	Clients         map[*websocket.Conn]bool
	mu              sync.Mutex
	Broadcast       chan message.FinalMessage
	RepliesBroadcast chan message.FinalMessageReply
}

// NewBroadcaster initializes a new Broadcaster instance.
func NewBroadcaster() *Broadcast {
	return &Broadcast{
		Clients:         make(map[*websocket.Conn]bool),
		Broadcast:       make(chan message.FinalMessage),
		RepliesBroadcast: make(chan message.FinalMessageReply),
	}
}

// RegisterClient adds a new WebSocket client.
func (b *Broadcast) RegisterClient(client *websocket.Conn) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Clients[client] = true
}

// RemoveClient removes a WebSocket client and closes the connection.
func (b *Broadcast) RemoveClient(client *websocket.Conn) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, exists := b.Clients[client]; exists {
		client.Close()
		delete(b.Clients, client)
	}
}

// HandleMessages listens for messages and broadcasts them to all clients.
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

// BroadcastMessage sends a message to all clients.
func (b *Broadcast) BroadcastMessage(msg message.FinalMessage) {
	select {
	case b.Broadcast <- msg:
	default:
		log.Println("Broadcast channel full, dropping message")
	}
}

// BroadcastReplyMessage sends a reply message to all clients.
func (b *Broadcast) BroadcastReplyMessage(msg message.FinalMessageReply) {
	select {
	case b.RepliesBroadcast <- msg:
	default:
		log.Println("Reply broadcast channel full, dropping message")
	}
}
