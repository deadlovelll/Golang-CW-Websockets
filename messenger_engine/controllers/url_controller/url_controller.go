import (
	"context"
	"encoding/json"
	"fmt"
	"messenger_engine/models/message"
	"net/http"
)

// MessengerRepository defines the interface for database operations
type MessengerRepository interface {
	GetUserChats(ctx context.Context, userID int) ([]message.ChatMessage, error)
}

// PresignedURLService handles fetching user avatar URLs
type PresignedURLService struct {
	Client *http.Client
}

// GetPresignedURL fetches the presigned avatar URL for a given user
func (s *PresignedURLService) GetPresignedURL(ctx context.Context, userID int) (string, error) {
	url := fmt.Sprintf("http://127.0.0.1:8165?user_id=%d", userID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		Status       string `json:"STATUS"`
		PresignedURL string `json:"PRESIGNED_URL"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error decoding JSON: %w", err)
	}

	return response.PresignedURL, nil
}