package main

import (
	"log"
	"net/http"

	"hashtags_search/controllers/base_controller"
	"hashtags_search/controllers/hashtag_controller"
	"hashtags_search/modules/database/database_pool"
	"hashtags_search/server"
	"hashtags_search/handlers/websocket_handler"
)

// main initializes the database, sets up controllers, registers routes, and starts the HTTP server.
func main() {
	// Initialize the database pool.
	dbPool := &databasepool.DatabasePoolController{}
	dbPool.StartupEvent()
	defer dbPool.ShutdownEvent()

	// Initialize controllers for handling database interactions.
	baseCtrl := basecontroller.BaseController{Database: dbPool.GetDb()}
	getDbCtrl := hashtagcontroller.HashtagController{BaseController: baseCtrl}

	// Create WebSocket handler for real-time communication.
	wsHandler := websockethandler.NewWebSocketHandler(&getDbCtrl)

	// Register HTTP routes.
	http.Handle("/", wsHandler)

	// Define server address and start the server.
	addr := "localhost:8380"
	log.Printf("Starting server on http://%s", addr)
	server.StartServer(addr, http.DefaultServeMux)
}
