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

func main() {
	// Initialize the database pool.
	dbPool := database() // Assuming you have a constructor.
	dbPool.StartupEvent()
	defer dbPool.ShutdownEvent()

	// Initialize controllers.
	baseCtrl := basecontroller.BaseController{Database: dbPool.GetDb()}
	dbCtrl := placecontroller.PlaceController{BaseController: &baseCtrl}

	// Set up the WebSocket handler.
	wsHandler := websockethandler.NewWebSocketHandler(&dbCtrl)

	// Register the handler and start the HTTP server.
	http.Handle("/", wsHandler)
	server := &http.Server{Addr: "localhost:8285"}

	// Start the HTTP server in a separate goroutine.
	go func() {
		log.Println("Server starting on http://localhost:8285")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for termination signals.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Gracefully shut down the server.
	if err := server.Close(); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}
	log.Println("Server gracefully stopped.")
}
