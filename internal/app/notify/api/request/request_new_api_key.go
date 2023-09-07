package request

// APINewKeyRequest specifies the necessary fields for creating a new APIKey.
type APINewKeyRequest struct {
	// Message is the message to be sent.
	Role uint8 `json:"Role"`
}

func (req *APINewKeyRequest) IsValid() (bool, error) {
	return true, nil
}
