package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"messenger_engine/models/presigned_url"
	"messenger_engine/controllers/url_controller"
)

// Test successful fetching of a presigned URL
func TestFetch_Success(t *testing.T) {
	// Expected response
	expectedResponse := presignedurl.PresignedUrlResponse{
		PresignedURL: "https://example.com/presigned-url",
	}

	// Create a test server that returns the expected response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure the request method is GET
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	}))
	defer server.Close()

	// Create an instance of HttpPresignedUrlFetcher
	fetcher := urlcontroller.HttpPresignedUrlFetcher{}

	// Call Fetch with the test server URL
	resp, err := fetcher.Fetch(server.URL)
	if err != nil {
		t.Fatalf("Fetch() returned an unexpected error: %v", err)
	}

	// Validate response
	if resp.PresignedURL != expectedResponse.PresignedURL {
		t.Errorf("expected PresignedURL %s, got %s", expectedResponse.PresignedURL, resp.PresignedURL)
	}
}

// Test fetching when the server returns invalid JSON
func TestFetch_InvalidJSON(t *testing.T) {
	// Create a test server that returns malformed JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"invalid_json": `) // Malformed JSON
	}))
	defer server.Close()

	// Create an instance of HttpPresignedUrlFetcher
	fetcher := urlcontroller.HttpPresignedUrlFetcher{}

	// Call Fetch and expect an error
	_, err := fetcher.Fetch(server.URL)
	if err == nil {
		t.Fatalf("Fetch() expected an error due to invalid JSON, but got none")
	}
}

// Test fetching when the request fails (simulating a network error)
func TestFetch_RequestFailure(t *testing.T) {
	// Use an invalid URL to simulate a network failure
	invalidURL := "http://invalid.invalid"

	// Create an instance of HttpPresignedUrlFetcher
	fetcher := urlcontroller.HttpPresignedUrlFetcher{}

	// Call Fetch and expect an error
	_, err := fetcher.Fetch(invalidURL)
	if err == nil {
		t.Fatalf("Fetch() expected an error due to network failure, but got none")
	}
}
