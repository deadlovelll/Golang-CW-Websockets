package presignedurl

// PresignedUrlResponse represents the response structure for a presigned URL request.
//
// Fields:
//   - Status: The status of the request (e.g., "success", "error").
//   - PresignedURL: The generated presigned URL for accessing a resource.
type PresignedUrlResponse struct {
	Status       string `json:"STATUS"`
	PresignedURL string `json:"PRESIGNED_URL"`
}
