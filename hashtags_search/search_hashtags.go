package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	getdatabasecontroller "hashtags_search/get_database_controller"
	basecontroller "hashtags_search/modules/base_controller"
	"hashtags_search/modules/database"

	"github.com/gorilla/websocket"
)

// WebSocketHandler handles WebSocket connections and processes incoming messages.
type WebSocketHandler struct {
	upgrader websocket.Upgrader
	dbCtrl   *getdatabasecontroller.GetDatabaseController
}

// Message represents the JSON structure for incoming WebSocket messages.
type Message struct {
	Query string `json:"query"`
}

// ServeHTTP upgrades the HTTP connection to a WebSocket, then continuously listens
// for and processes incoming messages. It uses the get-database controller to process
// queries and sends back JSON responses.
func (wsh *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Allow all origins (customize this for production use).
	wsh.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// Upgrade the HTTP connection to a WebSocket.
	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Process incoming messages in an infinite loop.
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading WebSocket message: %v", err)
			break
		}

		// Parse the JSON message.
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error parsing JSON message: %v", err)
			continue
		}

		// Use the query string as the hashtag.
		hashtag := msg.Query

		// Fetch data using the get-database controller.
		jsonData, err := wsh.dbCtrl.GetHashtags(hashtag)
		if err != nil {
			log.Printf("Error fetching hashtags for query '%s': %v", hashtag, err)
			return
		}

		// Log the processed query and response.
		log.Printf("Received query: %s; Responding with: %s", hashtag, string(jsonData))

		// Write the JSON response back to the client.
		if err := conn.WriteMessage(messageType, jsonData); err != nil {
			log.Printf("Error sending WebSocket message: %v", err)
			break
		}
	}
}

func main() {
	// Initialize the database pool.
	dbPool := &database.DatabasePoolController{}
	dbPool.StartupEvent()
	// Ensure graceful shutdown of the DB pool.
	defer dbPool.ShutdownEvent()

	// Initialize the base controller and get-database controller.
	baseCtrl := basecontroller.BaseController{Database: dbPool.GetDb()}
	getDbCtrl := getdatabasecontroller.GetDatabaseController{BaseController: &baseCtrl}

	// Create a WebSocketHandler with the get-database controller.
	wsHandler := &WebSocketHandler{
		upgrader: websocket.Upgrader{},
		dbCtrl:   &getDbCtrl,
	}

	// Register the WebSocketHandler and create an HTTP server.
	http.Handle("/", wsHandler)
	addr := "localhost:8380"
	log.Printf("Starting server on http://%s", addr)

	server := &http.Server{
		Addr: addr,
	}

	// Run the HTTP server in a separate goroutine.
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Set up channel to listen for interrupt or termination signals.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received.")

	// Create a context with timeout for graceful shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server gracefully stopped.")
}
