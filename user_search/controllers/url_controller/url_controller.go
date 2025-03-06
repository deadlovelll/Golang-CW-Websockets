package urlcontroller

import (
	"encoding/json"
	"fmt"
	"net/http"

	Response "user_search/models/presigned_url"
)


type UrlResponse struct {
	Status       string `json:"STATUS"`
	PresignedURL string `json:"PRESIGNED_URL,omitempty"`
}

// HttpPresignedUrlFetcher is responsible for fetching presigned URLs
// using the net/http package.
type HttpPresignedUrlFetcher struct{}

// Fetch makes an HTTP GET request to retrieve a presigned URL and decodes the response.
//
// Parameters:
//   - url: The URL from which to fetch the presigned URL response.
//
// Returns:
//   - A pointer to Response.PresignedUrlResponse containing the retrieved presigned URL.
//   - An error if the request fails or the response cannot be decoded.
func (h *HttpPresignedUrlFetcher) Fetch(url string) (*Response.PresignedUrlResponse, error) {
	// Send HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Decode JSON response
	var response Response.PresignedUrlResponse
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return &response, nil
}
