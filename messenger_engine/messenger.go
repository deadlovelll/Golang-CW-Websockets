package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"

	"messenger_engine/modules/database"
	BaseController "messenger_engine/controllers/base_controller"
	WebSocketController "messenger_engine/controllers/websocket_controller"
	MessageCtrl "messenger_engine/controllers/message_controller"
	MessageTypeController "messenger_engine/controllers/message_type_controller"
)

// MainMessageController combines processing of various WebSocket message types.
type MainMessageController struct {
	upgrader              websocket.Upgrader
	WsController          WebSocketController.WebsocketRepository
	MessageTypeController MessageTypeController.MessageTypeRepository
	Repo                  MessageCtrl.MessageRepository
}

// ServeHTTP upgrades the HTTP connection to a WebSocket, registers the client, and handles incoming messages.
func (mc *MainMessageController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Allow all origins (for development; restrict in production)
	mc.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := mc.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error upgrading to WebSocket: %s", err), http.StatusInternalServerError)
		return
	}
	defer ws.Close()

	// Register the client using the controller's method (do not modify the clients map directly)
	mc.WsController.AddClient(ws)

	// Main loop for reading messages
	for {
		var msg map[string]interface{}
		if err := ws.ReadJSON(&msg); err != nil {
			log.Printf("Error reading message: %s", err)
			mc.WsController.RemoveClient(ws)
			ws.WriteJSON(map[string]string{"error": fmt.Sprintf("Error reading message: %s", err)})
			break
		}

		// Dispatch the message based on its type
		switch msg["type"] {
		case "initial":
			mc.MessageTypeController.HandleInitialMessage(ws, msg)
		case "message":
			mc.MessageTypeController.HandleNewMessage(ws, msg)
		case "message_reply":
			mc.MessageTypeController.HandleMessageReply(ws, msg)
		default:
			ws.WriteJSON(map[string]string{"error": "Invalid message type"})
		}
	}
}

func main() {
	// Initialize the database pool
	dbPool := &database.DatabasePoolController{}
	if err := dbPool.StartupEvent(); err != nil {
		log.Fatalf("Failed to start database: %v", err)
	}
	defer func() {
		if err := dbPool.ShutdownEvent(); err != nil {
			log.Printf("Error during DB shutdown: %v", err)
		}
	}()

	// Initialize the base controller with the database connection
	baseCtrl := BaseController.BaseController{Database: dbPool.GetDb()}

	// Initialize controllers for retrieving and sending messages
	getCtrl := GetDatabaseController.GetMessengerController{BaseController: &baseCtrl}
	postCtrl := PostMessengerController.MakeMessagesController{BaseController: &baseCtrl}

	// Initialize WebSocket handlers.
	// Assume ChatsHandler implements WebsocketRepository.
	wsHandler := &ChatsHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		GMContrl: &getCtrl,
	}

	// Assume ChatMessageHandler implements MessageTypeRepository.
	chatMsgHandler := &ChatMessageHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		MMC: &postCtrl,
	}

	// Initialize the main message controller that combines message type handling.
	mainMsgCtrl := &MainMessageController{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		WsController:          wsHandler,      // ChatsHandler satisfies WebsocketRepository
		MessageTypeController: chatMsgHandler, // ChatMessageHandler satisfies MessageTypeRepository
		Repo:                  postCtrl,       // Use the message repository from postCtrl
	}

	// Register routes
	http.Handle("/chats", wsHandler)
	http.Handle("/chat", mainMsgCtrl)

	// Start message handling in the background
	go mainMsgCtrl.WsController.HandleMessages()

	// Configure and start the HTTP server
	server := &http.Server{Addr: "localhost:8440"}
	go func() {
		log.Printf("Server starting on http://%s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for termination signals (SIGINT, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Perform graceful shutdown with a 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped.")
}