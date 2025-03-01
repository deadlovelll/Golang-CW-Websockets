package broadcastcontroller

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"

	"messenger_engine/models/message"
	"messenger_engine/controllers/message_controller"
)

// Broadcaster manages WebSocket clients and message broadcasting.
type Broadcaster struct {
	clients         map[*websocket.Conn]bool
	mu              sync.Mutex
	broadcast       chan message.FinalMessage
	repliesBroadcast chan message.FinalMessageReply
}

// NewBroadcaster initializes a new Broadcaster instance.
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients:         make(map[*websocket.Conn]bool),
		broadcast:       make(chan message.FinalMessage),
		repliesBroadcast: make(chan message.FinalMessageReply),
	}
}

// RegisterClient adds a new WebSocket client.
func (b *Broadcaster) RegisterClient(client *websocket.Conn) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clients[client] = true
}

// RemoveClient removes a WebSocket client and closes the connection.
func (b *Broadcaster) RemoveClient(client *websocket.Conn) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, exists := b.clients[client]; exists {
		client.Close()
		delete(b.clients, client)
	}
}

// HandleMessages listens for messages and broadcasts them to all clients.
func (b *Broadcaster) HandleMessages(mmc *messagecontroller.MessageController) {
	for msg := range b.broadcast {
		b.mu.Lock()
		for client := range b.clients {
			if err := client.WriteJSON(msg); err != nil {
				log.Printf("Error sending message to client: %v", err)
				b.RemoveClient(client)
			}
		}
		b.mu.Unlock()
	}
}

// BroadcastMessage sends a message to all clients.
func (b *Broadcaster) BroadcastMessage(msg message.FinalMessage) {
	select {
	case b.broadcast <- msg:
	default:
		log.Println("Broadcast channel full, dropping message")
	}
}

// BroadcastReplyMessage sends a reply message to all clients.
func (b *Broadcaster) BroadcastReplyMessage(msg message.FinalMessageReply) {
	select {
	case b.repliesBroadcast <- msg:
	default:
		log.Println("Reply broadcast channel full, dropping message")
	}
}
