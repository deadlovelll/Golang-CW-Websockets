package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"places_search/models"
)

// GetPresignedURL makes an HTTP GET request to the specified URL and parses the response.
// It expects the response body to contain JSON data that matches the structure of models.PlaceResponse.
//
// Parameters:
//   - url: A string representing the endpoint to fetch the presigned URL.
//
// Returns:
//   - *models.PlaceResponse: A pointer to the parsed PlaceResponse struct containing the response data.
//   - error: An error if the request fails or JSON decoding encounters an issue.
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
