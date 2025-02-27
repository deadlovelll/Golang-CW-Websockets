package urlcontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// URLFetcher defines the behavior for fetching a presigned URL.
type URLFetcher interface {
	GetPresignedURL(ctx context.Context, userID int) (string, error)
}

// HTTPClient defines an interface for making HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// URLResponseParser defines an interface for parsing presigned URL responses.
type URLResponseParser interface {
	ParseResponse(resp *http.Response) (string, error)
}

// PresignedURLFetcher fetches a presigned URL using an HTTP client.
type PresignedURLFetcher struct {
	Client HTTPClient
	Parser URLResponseParser
}

// NewPresignedURLFetcher creates a new instance of PresignedURLFetcher.
func NewPresignedURLFetcher(client HTTPClient, parser URLResponseParser) *PresignedURLFetcher {
	return &PresignedURLFetcher{Client: client, Parser: parser}
}

// Base URL format for fetching the presigned URL.
const presignedURLFormat = "http://127.0.0.1:8165?user_id=%d"

// GetPresignedURL fetches the presigned avatar URL for a given user ID.
func (p *PresignedURLFetcher) GetPresignedURL(ctx context.Context, userID int) (string, error) {
	// Construct the URL with the user ID.
	url := fmt.Sprintf(presignedURLFormat, userID)

	// Create a new HTTP request with context.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Send the request using the HTTP client.
	resp, err := p.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Parse the response using the parser.
	return p.Parser.ParseResponse(resp)
}

// JSONPresignedURLParser parses JSON responses containing presigned URLs.
type JSONPresignedURLParser struct{}

// ParseResponse extracts the presigned URL from the JSON response.
func (j *JSONPresignedURLParser) ParseResponse(resp *http.Response) (string, error) {
	// Check if the response status is OK (200).
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}

	// Define a struct to match the expected JSON format.
	var res struct {
		Status       string `json:"STATUS"`
		PresignedURL string `json:"PRESIGNED_URL"`
	}

	// Decode the JSON response.
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("error decoding JSON: %w", err)
	}

	// Return the extracted presigned URL.
	return res.PresignedURL, nil
}