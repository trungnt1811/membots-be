package response

// APIAllWhitelistedIPsResponse is a response for GetAllWhitelistedIPs requests.
type APIAllWhitelistedIPsResponse struct {
	Result map[string]string `json:"result"`
}
