package tests

import (
	"context"
	"net/http"
	"testing"
	"time"

	"hashtags_search/server"
)

// TestStartServer verifies that the server starts and responds correctly.
func TestStartServer(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	serverAddr := "127.0.0.1:8081"

	// Start server in a separate goroutine.
	go server.StartServer(serverAddr, handler)

	// Allow server time to start.
	time.Sleep(1 * time.Second)

	// Send a test request to the server.
	resp, err := http.Get("http://" + serverAddr + "/test")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

// TestServerShutdown verifies that the server shuts down gracefully.
func TestServerShutdown(t *testing.T) {
	handler := http.NewServeMux()
	serverAddr := "127.0.0.1:8082"

	httpServer := &http.Server{
		Addr:    serverAddr,
		Handler: handler,
	}

	// Start server in a separate goroutine.
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Fatalf("Server error: %v", err)
		}
	}()

	// Allow server time to start.
	time.Sleep(1 * time.Second)

	// Create a shutdown context.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		t.Fatalf("Server shutdown error: %v", err)
	}
}
