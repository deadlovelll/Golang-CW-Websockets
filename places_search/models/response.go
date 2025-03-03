package models

// PlaceResponse represents the JSON response for a presigned URL.
type PlaceResponse struct {
	Status       string   `json:"STATUS"`
	PresignedURL []string `json:"PRESIGNED_URL"`
}
