package tests

import (
	"net/http"
	"messenger_engine/controllers/broadcast_controller"
	"messenger_engine/models/message"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
)

// TestRegisterClient verifies that a WebSocket client can successfully connect 
// and be registered in the broadcaster's client list.
func TestRegisterClient(t *testing.T) {
	broadcaster := broadcastcontroller.NewBroadcaster()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _ := websocket.Upgrade(w, r, nil, 1024, 1024)
		broadcaster.RegisterClient(conn)
	}))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]
	client, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect WebSocket: %v", err)
	}
	defer client.Close()

	if _, exists := broadcaster.Clients[client]; !exists {
		t.Errorf("Client was not registered correctly")
	}
}

// TestBroadcastMessage verifies that a message sent via the broadcaster is 
// received by the connected WebSocket clients.
func TestBroadcastMessage(t *testing.T) {
	broadcaster := broadcastcontroller.NewBroadcaster()
	var wg sync.WaitGroup

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _ := websocket.Upgrade(w, r, nil, 1024, 1024)
		broadcaster.RegisterClient(conn)
		wg.Add(1)
		go func() {
			defer wg.Done()
			var receivedMsg message.FinalMessage
			err := conn.ReadJSON(&receivedMsg)
			if err != nil {
				t.Errorf("Failed to read message: %v", err)
			}
		}()
	}))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]
	client, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect WebSocket: %v", err)
	}
	defer client.Close()

	msg := message.FinalMessage{
		Type: "message",
		Message: message.Message{ 
			Message: "Hello, World!",
		},
	}
	broadcaster.BroadcastMessage(msg)
	wg.Wait()
}

// TestRemoveClient verifies that a WebSocket client is properly removed from the 
// broadcaster's client list and its connection is closed.
func TestRemoveClient(t *testing.T) {
	broadcaster := broadcastcontroller.NewBroadcaster()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _ := websocket.Upgrade(w, r, nil, 1024, 1024)
		broadcaster.RegisterClient(conn)
		broadcaster.RemoveClient(conn)
	}))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]
	client, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect WebSocket: %v", err)
	}
	defer client.Close()

	if _, exists := broadcaster.Clients[client]; exists {
		t.Errorf("Client was not removed correctly")
	}
}
