package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"places_search/models"
)

// GetPresignedURL makes an HTTP GET request to the provided URL and decodes the JSON response.
func GetPresignedURL(url string) (*models.PlaceResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	var response models.PlaceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}
	return &response, nil
}
