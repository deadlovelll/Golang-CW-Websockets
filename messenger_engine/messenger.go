package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"messenger_engine/modules/database"

	BaseController "messenger_engine/controllers/base_controller"
	ChatController "messenger_engine/controllers/chat_controller"
	MessageController "messenger_engine/controllers/message_controller"
	Broadcast "messenger_engine/controllers/broadcast_controller"

	"github.com/gorilla/websocket"

	ChatMessageHandler "messenger_engine/controllers/websocket_controller/handlers/chat_message_handler"
	ChatHandler "messenger_engine/controllers/websocket_controller/handlers/chat_handler"
)

func main() {
	// Initialize the database pool.
	dbPool := &database.DatabasePoolController{}
	dbPool.StartupEvent()

	// Create a base controller instance.
	baseCtrl := BaseController.BaseController{Database: dbPool.GetDb()}

	// Initialize chat and message controllers.
	chatCtrl := ChatController.ChatController{BaseController: &baseCtrl}
	messageCtrl := MessageController.MessageController{BaseController: &baseCtrl}

	// Create WebSocket handlers with dependency injection.
	wsHandler := ChatHandler.NewChatsHandler(websocket.Upgrader{}, &chatCtrl)
	chatMsgHandler := ChatMessageHandler.NewChatMessageHandler(websocket.Upgrader{}, &messageCtrl)

	// Setup HTTP routes.
	http.Handle("/chats", wsHandler)
	http.Handle("/chat", chatMsgHandler)

	log.Println("Starting server on http://localhost:8440")

	// Start the broadcast routine.
	go Broadcast.HandleMessages(&messageCtrl)

	// Start the HTTP server.
	server := &http.Server{Addr: "localhost:8440"}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting server: %v", err)
		}
	}()

	// Wait for termination signals.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	dbPool.ShutdownEvent()

	if err := server.Close(); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	log.Println("Server successfully shut down.")
}
