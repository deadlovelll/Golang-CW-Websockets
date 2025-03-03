package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"places_search/handlers/websocket_handler"
	"places_search/controllers/base_controller"
	"places_search/modules/database/database"
	"places_search/controllers/place_controller"
)

// main is the entry point of the application.
// It initializes the database, controllers, and WebSocket handler,
// then starts the HTTP server and listens for termination signals.
func main() {
	// Initialize the database pool (singleton instance).
	dbPool := database.GetDatabaseInstance()

	// Initialize controllers with the database instance.
	baseCtrl := basecontroller.BaseController{Database: dbPool}
	dbCtrl := placecontroller.PlaceController{BaseController: &baseCtrl}

	// Set up the WebSocket handler with the PlaceController.
	wsHandler := websockethandler.NewWebSocketHandler(&dbCtrl)

	// Register the WebSocket handler and start the HTTP server.
	http.Handle("/", wsHandler)
	server := &http.Server{Addr: "localhost:8285"}

	// Start the HTTP server in a separate goroutine.
	go func() {
		log.Println("Server starting on http://localhost:8285")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Listen for OS termination signals (SIGINT, SIGTERM) to handle graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Gracefully shut down the server when a termination signal is received.
	if err := server.Close(); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}
	log.Println("Server gracefully stopped.")
}