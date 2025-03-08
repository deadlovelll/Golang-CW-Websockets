package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"user_search/controllers/url_controller"
)

// TestFetch_Success tests the successful scenario where the HTTP GET request returns valid JSON data.
// It sets up an httptest server that responds with a JSON payload containing a "STATUS" of "SUCCESS"
// and a valid "PRESIGNED_URL". The function then verifies that the Fetch method correctly decodes the JSON
// response and returns the expected values.
func TestFetch_Success(t *testing.T) {
	// Create an httptest server that returns a valid JSON response.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Return JSON matching the expected structure.
		fmt.Fprintln(w, `{"STATUS": "SUCCESS", "PRESIGNED_URL": "http://example.com"}`)
	}))
	defer ts.Close()

	fetcher := &urlcontroller.HttpPresignedUrlFetcher{}
	res, err := fetcher.Fetch(ts.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	// Check the response fields.
	if res.Status != "SUCCESS" {
		t.Errorf("Expected status 'SUCCESS', got '%s'", res.Status)
	}
	if res.PresignedURL != "http://example.com" {
		t.Errorf("Expected URL 'http://example.com', got '%s'", res.PresignedURL)
	}
}

// TestFetch_InvalidJSON tests the scenario where the HTTP GET request returns an invalid JSON response.
// An httptest server is set up to send malformed JSON data. The test verifies that the Fetch method
// detects the decoding error and returns an appropriate error message.
func TestFetch_InvalidJSON(t *testing.T) {
	// Create an httptest server that returns invalid JSON.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `invalid json`)
	}))
	defer ts.Close()

	fetcher := &urlcontroller.HttpPresignedUrlFetcher{}
	_, err := fetcher.Fetch(ts.URL)
	if err == nil {
		t.Errorf("Expected error decoding JSON, got nil")
	}
}

// TestFetch_RequestError tests the scenario where the HTTP GET request fails due to an invalid URL.
// This test forces a request error by providing an invalid URL, and then verifies that the Fetch method
// returns an error, ensuring that error handling is properly implemented.
func TestFetch_RequestError(t *testing.T) {
	// Use an invalid URL to force an HTTP request error.
	fetcher := &urlcontroller.HttpPresignedUrlFetcher{}
	_, err := fetcher.Fetch("http://127.0.0.1:0")
	if err == nil {
		t.Errorf("Expected error for an invalid URL, got nil")
	}
}
