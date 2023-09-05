package request

// APIKeyRequest holds an APIKey.
type APIKeyRequest struct {
	// Message is the message to be sent.
	APIKey string `json:"apiKey" binding:"required" validate:"required"`
}

func (req *APIKeyRequest) IsValid() (bool, error) {
	return true, nil
}
