package urlcontroller

import (
	"encoding/json"
	"fmt"
	"net/http"

	Response "messenger_engine/models/presigned_url"
)

// HttpPresignedUrlFetcher implements PresignedUrlFetcher using the net/http package.
type HttpPresignedUrlFetcher struct{}

// Fetch makes an HTTP GET request to retrieve a presigned URL.
func (h *HttpPresignedUrlFetcher) Fetch(url string) (*Response.PresignedUrlResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var response Response.PresignedUrlResponse
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return &response, nil
}
