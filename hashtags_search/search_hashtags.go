package main

import (
	"log"
	"net/http"

	"hashtags_search/controllers/base_controller"
	"hashtags_search/controllers/get_database_controller"
	"hashtags_search/modules/database/database_pool"
	"hashtags_search/server"
	"hashtags_search/handlers/websocket_handler"
)

func main() {
	// Initialize the database pool.
	dbPool := &databasepool.DatabasePoolController{}
	dbPool.StartupEvent()
	defer dbPool.ShutdownEvent()

	// Initialize controllers.
	baseCtrl := basecontroller.BaseController{Database: dbPool.GetDb()}
	getDbCtrl := getdatabasecontroller.GetDatabaseController{BaseController: &baseCtrl}

	// Create WebSocket handler.
	wsHandler := websocket.NewWebSocketHandler(&getDbCtrl)

	// Register routes.
	http.Handle("/", wsHandler)

	// Start the HTTP server.
	addr := "localhost:8380"
	log.Printf("Starting server on http://%s", addr)
	server.StartServer(addr, http.DefaultServeMux)
}
