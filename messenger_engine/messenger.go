package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"

	"messenger_engine/modules/database/database_pool"

	// Controllers
	"messenger_engine/controllers/base_controller"
	"messenger_engine/controllers/chat_controller"
	"messenger_engine/controllers/message_controller"
	"messenger_engine/controllers/broadcast_controller"

	// WebSocket Handlers
	"messenger_engine/controllers/websocket_controller/handlers/chat_message_handler"
	"messenger_engine/controllers/websocket_controller/handlers/chat_handler"

	"messenger_engine/utls/env"
)

const serverAddr = "localhost:8440"

// main is the entry point of the application. It initializes environment variables,
// database connections, controllers, WebSocket handlers, and starts the HTTP server.
func main() {
	// Load environment variables from .env file
	goenv.LoadEnv()

	// Initialize the database connection pool
	dbPool := initializeDatabase()
	defer dbPool.StartupEvent()

	// Initialize controllers
	baseCtrl := basecontroller.BaseController{Database: dbPool.GetDb()}
	chatCtrl := chatcontroller.ChatController{BaseController: &baseCtrl}
	messageCtrl := messagecontroller.MessageController{BaseController: &baseCtrl}
	broadcastCtrl := broadcastcontroller.NewBroadcaster()
	
	// Initialize WebSocket handlers
	wsHandler := chathandler.NewChatsHandler(websocket.Upgrader{}, &chatCtrl)
	chatMsgHandler := chatmessagehandler.NewChatMessageHandler(websocket.Upgrader{}, &messageCtrl, broadcastCtrl)

	// Configure HTTP routes
	mux := http.NewServeMux()
	mux.Handle("/chats", wsHandler)
	mux.Handle("/chat", chatMsgHandler)

	// Start message broadcasting routine
	go broadcastcontroller.NewBroadcaster().HandleMessages(&messageCtrl)

	// Start HTTP server with graceful shutdown handling
	startServer(mux)
}

// initializeDatabase sets up and returns a new database pool controller.
func initializeDatabase() *databasepool.DatabasePoolController {
	dbPool := &databasepool.DatabasePoolController{}
	dbPool.StartupEvent()
	return dbPool
}

// startServer initializes and starts the HTTP server in a separate goroutine.
// It also triggers the graceful shutdown handling mechanism.
//
// Parameters:
//   - handler: The HTTP handler (mux router) to handle incoming requests.
func startServer(handler http.Handler) {
	server := &http.Server{
		Addr:    serverAddr,
		Handler: handler,
	}

	// Run server in a separate goroutine
	go func() {
		log.Printf("Server started on http://%s", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Handle graceful shutdown
	waitForShutdown(server)
}

// waitForShutdown listens for termination signals and gracefully shuts down the HTTP server.
//
// Parameters:
//   - server: The HTTP server instance to be gracefully shut down.
func waitForShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server successfully shut down.")
}