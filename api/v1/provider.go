package v1

type Provider struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	EndPoint string   `json:"endPoint"`
	APIKey   string   `json:"apiKey"`
	Weight   int      `json:"weight"`
	Models   []string `json:"models"`
}

type CreateProviderRequest struct {
	Name     string   `json:"name"`
	Type     string   `json:"type" binding:"required;oneof=openai,azure"`
	EndPoint string   `json:"end_point" binding:"required"`
	APIKey   string   `json:"api_key" binding:"required"`
	Weight   int      `json:"weight" default:"10"`
	Models   []string `json:"models"`
}
