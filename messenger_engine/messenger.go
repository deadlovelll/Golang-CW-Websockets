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

func main() {

	// Loading the .env file vars
	goenv.LoadEnv()

	// Initialize database
	dbPool := initializeDatabase()
	defer dbPool.StartupEvent()

	// Initialize controllers
	baseCtrl := basecontroller.BaseController{Database: dbPool.GetDb()}
	chatCtrl := chatcontroller.ChatController{BaseController: &baseCtrl}
	messageCtrl := messagecontroller.MessageController{BaseController: &baseCtrl}
	broadcastCtrl := broadcastcontroller.Broadcast{}

	// Initialize WebSocket handlers
	wsHandler := chathandler.NewChatsHandler(websocket.Upgrader{}, &chatCtrl)
	chatMsgHandler := chatmessagehandler.NewChatMessageHandler(websocket.Upgrader{}, &messageCtrl, &broadcastCtrl)

	// Configure HTTP routes
	mux := http.NewServeMux()
	mux.Handle("/chats", wsHandler)
	mux.Handle("/chat", chatMsgHandler)

	// Start broadcast routine
	go broadcastcontroller.NewBroadcaster().HandleMessages(&messageCtrl)

	// Start server with graceful shutdown
	startServer(mux)
}

// initializeDatabase sets up the database pool
func initializeDatabase() *databasepool.DatabasePoolController {
	dbPool := &databasepool.DatabasePoolController{}
	dbPool.StartupEvent()
	return dbPool
}

// startServer starts the HTTP server and handles graceful shutdown
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

// waitForShutdown handles termination signals and graceful server shutdown
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
