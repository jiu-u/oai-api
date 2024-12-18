package v1

// doc https://openai.apifox.cn/api-55352401

type CompletionsRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Echo   bool   `json:"echo,omitempty"`

	Stream           bool    `json:"stream,omitempty"`
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
	Logprobs         any     `json:"logprobs,omitempty"`
	MaxTokens        int64   `json:"max_tokens,omitempty"`
	Temperature      float64 `json:"temperature,omitempty"`
	BestOf           int64   `json:"best_of,omitempty"`
	LogitBias        any     `json:"logit_bias,omitempty"`
	N                any     `json:"n,omitempty"`
	PresencePenalty  float64 `json:"presence_penalty,omitempty"`
	Seed             any     `json:"seed,omitempty"`
	Stop             any     `json:"stop,omitempty"`
	Suffix           string  `json:"suffix,omitempty"`
	TopP             float64 `json:"top_p,omitempty"`
	User             string  `json:"user,omitempty"`
}

type (
	CompletionsResp struct {
		ID                string              `json:"id"`
		Object            string              `json:"object"`
		Created           int64               `json:"created"`
		Model             string              `json:"model"`
		Choices           []CompletionsChoice `json:"choices"`
		Usage             Usage               `json:"usage,omitempty"`
		SystemFingerprint string              `json:"systemFingerprint,omitempty"`
	}
	CompletionsChoice struct {
		Text         string `json:"text"`
		Index        int    `json:"index"`
		Logprobs     any    `json:"logprobs,omitempty"`
		FinishReason string `json:"finish_reason,omitempty"`
	}
)
