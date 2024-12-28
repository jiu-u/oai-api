package v1

type CreateApiKeyRequest struct {
	UserId string `json:"userId"`
}

type CreateApiKeyResponse struct {
	ApiKey string `json:"apiKey"`
}

type ResetApiKeyRequest struct {
	UserId string `json:"userId"`
}

type ResetApiKeyResponse struct {
	ApiKey string `json:"apiKey"`
}
