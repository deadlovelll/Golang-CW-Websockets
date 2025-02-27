package websocketcontroller

import (
	"log"
	"sync"

	Messages "messenger_engine/models/message"

	"github.com/gorilla/websocket"
)

type WebsocketRepository interface {
	HandleMessages()
	AddClient(client *websocket.Conn)
	RemoveClient(client *websocket.Conn)
	GetClients() map[*websocket.Conn]bool
}

// MessageBroadcaster manages WebSocket clients and message broadcasting.
type MessageBroadcaster struct {
	mu              sync.RWMutex
	clients         map[*websocket.Conn]bool
	broadcast       chan Messages.FinalMessage
	repliesBroadcast chan Messages.FinalMessageReply
}

// AddClient registers a new WebSocket connection.
func (mb *MessageBroadcaster) AddClient(client *websocket.Conn) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	mb.clients[client] = true
}

// RemoveClient safely removes a client from the pool.
func (mb *MessageBroadcaster) RemoveClient(client *websocket.Conn) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	client.Close()
	delete(mb.clients, client)
}

func (mb *MessageBroadcaster) GetClients() map[*websocket.Conn]bool {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	clientsCopy := make(map[*websocket.Conn]bool)
	for client, active := range mb.clients {
		clientsCopy[client] = active
	}
	return clientsCopy
}

// HandleMessages listens for messages and broadcasts them to all clients.
func (mb *MessageBroadcaster) HandleMessages() {
	for msg := range mb.broadcast {
		mb.mu.RLock()
		for client := range mb.clients {
			if err := client.WriteJSON(msg); err != nil {
				log.Printf("Error sending message: %v", err)
				mb.mu.RUnlock()
				mb.RemoveClient(client)
				mb.mu.RLock()
			}
		}
		mb.mu.RUnlock()
	}
}
