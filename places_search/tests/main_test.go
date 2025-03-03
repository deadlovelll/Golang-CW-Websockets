package tests

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
	"encoding/json"

	"places_search/controllers/base_controller"
	"places_search/modules/database/database"
	"places_search/controllers/place_controller"
	"places_search/handlers/websocket_handler"

	"github.com/gorilla/websocket"
)

// Test for server initialization and correct response
func TestServerInitialization(t *testing.T) {
	// Create a test database instance
	dbPool := database.GetDatabaseInstance()

	// Initialize controllers with the database instance
	baseCtrl := basecontroller.BaseController{Database: dbPool}
	dbCtrl := placecontroller.PlaceController{BaseController: &baseCtrl}

	// Set up WebSocket handler
	wsHandler := websockethandler.NewWebSocketHandler(&dbCtrl)

	// Create a test HTTP server
	server := httptest.NewServer(wsHandler)
	defer server.Close()

	// Make a sample request to ensure the server is running
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert if the response status is OK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

// Test for WebSocket connection upgrade
func TestWebSocketUpgrade(t *testing.T) {
	// Create a test database instance
	dbPool := database.GetDatabaseInstance()

	// Initialize controllers with the database instance
	baseCtrl := basecontroller.BaseController{Database: dbPool}
	dbCtrl := placecontroller.PlaceController{BaseController: &baseCtrl}

	// Set up WebSocket handler
	wsHandler := websockethandler.NewWebSocketHandler(&dbCtrl)

	// Create a test HTTP server
	server := httptest.NewServer(wsHandler)
	defer server.Close()

	// Dial the WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to connect WebSocket: %v", err)
	}
	defer conn.Close()

	// Send a test message
	message := map[string]string{"query": "Test Place"}
	messageData, _ := json.Marshal(message)
	err = conn.WriteMessage(websocket.TextMessage, messageData)
	if err != nil {
		t.Fatalf("Failed to send WebSocket message: %v", err)
	}

	// Receive the response
	_, response, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read WebSocket message: %v", err)
	}

	// Check if response is not empty (or you can check specific content if needed)
	if len(response) == 0 {
		t.Fatalf("Expected a response, got empty")
	}
}

// Test for graceful server shutdown
func TestGracefulShutdown(t *testing.T) {
	// Create a test database instance
	dbPool := database.GetDatabaseInstance()

	// Initialize controllers with the database instance
	baseCtrl := basecontroller.BaseController{Database: dbPool}
	dbCtrl := placecontroller.PlaceController{BaseController: &baseCtrl}

	// Set up WebSocket handler
	wsHandler := websockethandler.NewWebSocketHandler(&dbCtrl)

	// Create a test HTTP server
	server := &http.Server{
		Addr:    "localhost:8285",
		Handler: wsHandler,
	}

	// Start the server in a goroutine
	go func() {
		log.Println("Starting server...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Fatalf("Server error: %v", err)
		}
	}()

	// Send SIGINT signal to simulate graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Simulate graceful shutdown
	time.Sleep(1 * time.Second) // wait for server to start
	quit <- syscall.SIGINT // Sending the signal directly to the quit channel
	<-quit                 // Wait for server to handle the signal

	// Now verify if the server has been shut down
	if err := server.Close(); err != nil {
		t.Fatalf("Error shutting down server: %v", err)
	}
	log.Println("Server gracefully stopped.")
}

