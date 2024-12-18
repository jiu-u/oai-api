package v1

type EmbeddingRequest struct {
	Model          string `json:"model"`
	Input          string `json:"input"`
	User           string `json:"user,omitempty"`
	EncodingFormat string `json:"encoding_format,omitempty"`
	Dimensions     int    `json:"dimensions,omitempty"`
}

type EmbeddingResponse struct {
	Object    string `json:"object"`
	Data      []Data `json:"data"`
	Model     string `json:"model"`
	Usage     Usage  `json:"usage"`
	Embedding []any  `json:"embedding,omitempty"`
}

type Data struct {
	ID        string `json:"id,omitempty"`
	Object    string `json:"object"`
	Index     int    `json:"index"`
	Embedding []any  `json:"embedding,omitempty"`
}
