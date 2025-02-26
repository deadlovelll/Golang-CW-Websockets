package urlcontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// URLFetcher defines the behavior of fetching a presigned URL.
type URLFetcher interface {
	GetPresignedURL(ctx context.Context, userID int) (string, error)
}

// PresignedURLFetcher implements URLFetcher using an HTTP client.
type PresignedURLFetcher struct {
	Client *http.Client
}

// GetPresignedURL fetches the presigned avatar URL for a given user.
func (p *PresignedURLFetcher) GetPresignedURL(ctx context.Context, userID int) (string, error) {
	url := fmt.Sprintf("http://127.0.0.1:8165?user_id=%d", userID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}

	var response struct {
		Status       string `json:"STATUS"`
		PresignedURL string `json:"PRESIGNED_URL"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error decoding JSON: %w", err)
	}

	return response.PresignedURL, nil
}
