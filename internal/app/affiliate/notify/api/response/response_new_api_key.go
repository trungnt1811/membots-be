package response

// APINewKeyResponse holds necessary information of response of an APIKeyRequest.
type APINewKeyResponse struct {
	Token       string `json:"id,omitempty"`
	ExpiredTime string `json:"expiredTime,omitempty"`
}
