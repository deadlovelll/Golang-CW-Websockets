package presignedurl

type PresignedUrlResponse struct {
	Status       string `json:"STATUS"`
	PresignedURL string `json:"PRESIGNED_URL"`
}