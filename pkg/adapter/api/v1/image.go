package v1

type (
	CreateImageRequest struct {
		Prompt         string `json:"prompt"`
		Model          string `json:"model,omitempty"`
		N              int    `json:"n,omitempty"`
		Quality        string `json:"quality,omitempty"`
		ResponseFormat string `json:"response_format,omitempty"`
		Size           string `json:"size,omitempty"`
		Style          string `json:"style,omitempty"` //  one of vivid or natural
		User           string `json:"user,omitempty"`
	}
	CreateImageResponse struct {
		Created int64       `json:"created"`
		Data    []ImageData `json:"data"`
	}
	ImageData struct {
		URL           string `json:"url,omitempty"`
		B64JSON       string `json:"b64_json,omitempty"`
		RevisedPrompt string `json:"revised_prompt,omitempty"`
	}
)

type (
	EditImageRequest struct {
		Image          []byte `form:"image"`
		Prompt         string `form:"prompt"`
		Mask           []byte `form:"mask,omitempty"`
		Model          string `form:"model,omitempty"`
		N              int    `form:"n,omitempty"`
		Size           string `form:"size,omitempty"`
		ResponseFormat string `form:"response_format,omitempty"`
		User           string `form:"user,omitempty"`
	}
)

type (
	CreateImageVariationRequest struct {
		Image          []byte `form:"image"`
		Model          string `form:"model,omitempty"`
		N              int    `form:"n,omitempty"`
		Size           string `form:"size,omitempty"`
		ResponseFormat string `form:"response_format,omitempty"`
		User           string `form:"user,omitempty"`
	}
)
