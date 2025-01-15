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
