package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// StartServer initializes and starts an HTTP server on the specified address with the given handler.
// It also handles graceful shutdown when termination signals (SIGINT, SIGTERM) are received.
//
// Parameters:
//   - addr: The address where the server will listen (e.g., ":8080").
//   - handler: The HTTP handler to process incoming requests.
func StartServer(addr string, handler http.Handler) {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Start the server in a separate goroutine to allow handling shutdown signals concurrently.
	go func() {
		log.Printf("Server started on http://%s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Create a channel to listen for OS signals indicating termination (SIGINT, SIGTERM).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Wait until a termination signal is received.
	log.Println("Shutdown signal received.")

	// Create a context with a timeout for graceful shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server.
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server gracefully stopped.")
}