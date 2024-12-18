package v1

// post https://api.openai.com/v1/audio/speech
// 返回文件
// doc https://platform.openai.com/docs/api-reference/audio/createSpeech

type SpeechRequest struct {
	Model          string  `json:"model"`
	Input          string  `json:"input"`
	Voice          string  `json:"voice"`
	ResponseFormat string  `json:"response_format,omitempty"`
	Speed          float64 `json:"speed,omitempty"`
	Text           string  `json:"text"`
}

// post https://api.openai.com/v1/audio/transcriptions
// doc https://platform.openai.com/docs/api-reference/audio/createTranscription

type TranscriptionRequest struct {
	File                   []byte  `form:"file"`
	Model                  string  `form:"model"`
	Language               string  `form:"language,omitempty"`
	Prompt                 string  `form:"prompt,omitempty"`
	ResponseFormat         string  `form:"response_format,omitempty"`
	Temperature            float64 `form:"temperature,omitempty"` // 温度
	TimestampGranularities any     `form:"timestamp_granularity,omitempty"`
}

type TransTextResponse struct {
	Text string `json:"text"`
}

type TransObjectResponse struct {
	Language string `json:"language"`
	Duration string `json:"duration"`
	Text     string `json:"text"`
	Words    []any  `json:"words,omitempty"`
	Segments []any  `json:"segments,omitempty"`
}

type TransWord struct {
	Word  string  `json:"word"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

type TransSegment struct {
	Id               int     `json:"id"`
	Seek             int     `json:"seek"`
	Start            float64 `json:"start"`
	End              float64 `json:"end"`
	Text             string  `json:"text"`
	Tokens           []any   `json:"tokens,omitempty"`
	Temperature      float64 `form:"temperature,omitempty"` // 温度
	AvgLogprobs      float64 `json:"avg_logprobs,omitempty"`
	CompressionRatio float64 `json:"compressionRatio,omitempty"`
	NoSpeechProb     float64 `json:"noSpeechProb,omitempty"`
}

// post https://api.openai.com/v1/audio/translations
// `TransTextResponse`
// doc https://platform.openai.com/docs/api-reference/audio/createTranslation

type TranslationRequest struct {
	file           []byte  `form:"file"`
	Model          string  `form:"model"`
	Prompt         string  `form:"prompt,omitempty"`
	ResponseFormat string  `form:"response_format,omitempty"`
	Temperature    float64 `form:"temperature,omitempty"` // 温度
}
