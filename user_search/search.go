package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"user_search/handlers/websocket_handler"
	"user_search/modules/database/database_pool"

	basecontroller "user_search/controllers/base_controller"
	"user_search/controllers/user_controller"
)

func main() {
	
	// Initialize the database pool
	dbPool := &databasepool.DatabasePoolController{}
	dbPool.StartupEvent()
	defer dbPool.ShutdownEvent()

	// Initialize controllers
	baseCtrl := basecontroller.NewBaseController(dbPool.GetDb())
	userCtrl := usercontroller.NewUserController(baseCtrl)

	// Setup the WebSocket handler
	wsHandler := websockethandler.NewWebSocketHandler(userCtrl)

	// Register the WebSocket handler at the root URL
	http.Handle("/", wsHandler)
	serverAddr := "localhost:8280"
	fmt.Printf("Server is running at http://%s\n", serverAddr)
	server := &http.Server{Addr: serverAddr}

	// Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup error: %v", err)
		}
	}()

	// Wait for termination signals (Ctrl+C, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Gracefully shut down the server
	if err := server.Close(); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}

	fmt.Println("Server has been gracefully terminated.")
}
