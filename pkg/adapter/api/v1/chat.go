package v1

type (
	ChatCompletionRequest struct {
		Model               string         `json:"model"`                 // 模型名称
		Messages            []Message      `json:"messages"`              // 聊天消息
		Temperature         float64        `json:"temperature,omitempty"` // 温度
		Stream              bool           `json:"stream"`
		SteamOptions        any            `json:"steam_options,omitempty"`     // 是否流式返回
		TopP                float64        `json:"top_p,omitempty"`             // top_p
		Stop                any            `json:"stop,omitempty"`              // 停止标志
		MaxTokens           int            `json:"max_tokens,omitempty"`        // 最大 token 数
		MaxCompletionTokens int            `json:"max_completion,omitempty"`    // 最大完成数
		PresencePenalty     float64        `json:"presence_penalty,omitempty"`  // 存在惩罚
		FrequencyPenalty    float64        `json:"frequency_penalty,omitempty"` // 频率惩罚
		LogitBias           map[string]any `json:"logit_bias,omitempty"`
		ResponseFormat      any            `json:"response_format,omitempty"`
		Store               bool           `json:"store,omitempty"`
		ReasoningEffect     string         `json:"reasoning_effect,omitempty"`
		MetaData            string         `json:"meta_data,omitempty"`
		Logprobs            any            `json:"logprobs,omitempty"`
		TopLogprobs         int            `json:"top_logprobs,omitempty"`
		N                   int            `json:"n,omitempty"`
		Prediction          any            `json:"prediction,omitempty"`
		Audio               any            `json:"audio,omitempty"`
		Modelities          any            `json:"modelities,omitempty"`
		Seed                int            `json:"seed,omitempty"`
		ServiceTier         string         `json:"service_tier,omitempty"`
		Tools               any            `json:"tools,omitempty"`
		ToolChoice          any            `json:"tool_choice,omitempty"`
		ParallelToolCalls   any            `json:"parallel_tool_calls,omitempty"`
		User                string         `json:"user,omitempty"`
		FunctionCall        any            `json:"function_call,omitempty"`
		Functions           any            `json:"functions,omitempty"`
	}

	Message struct {
		Role         string `json:"role"`
		Content      any    `json:"content"`
		Name         string `json:"name,omitempty"`
		ToolCallId   string `json:"tool_call_id,omitempty"`
		Refusal      any    `json:"refusal,omitempty"`
		Audio        any    `json:"audio,omitempty"`
		ToolCalls    any    `json:"tool_calls,omitempty"`
		FunctionCall any    `json:"function_call,omitempty"`
	}
)

type (
	ChatCompletionNoStreamResponse struct {
		ID      string   `json:"id"`
		Object  string   `json:"object"`
		Created int64    `json:"created"`
		Model   string   `json:"model"`
		Choices []Choice `json:"choices"`
		Usage   Usage    `json:"usage,omitempty"`
	}
	ChatCompletionStreamResponse struct {
		ID      string            `json:"id"`
		Object  string            `json:"object"`
		Created int64             `json:"created"`
		Model   string            `json:"model"`
		Choices []ChoiceWithDelta `json:"choices"`
		Usage   Usage             `json:"usage,omitempty"`
	}
	Choice struct {
		Index        int     `json:"index"`
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	}
	ChoiceWithDelta struct {
		Index        int    `json:"index"`
		Delta        Delta  `json:"delta"`
		FinishReason string `json:"finish_reason"`
	}
	Delta struct {
		Role    string `json:"role"`
		Content any    `json:"content"`
	}

	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}
)
