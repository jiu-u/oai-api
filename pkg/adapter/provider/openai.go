package provider

import (
	"bytes"
	"context"
	"github.com/bytedance/sonic"
	v1 "github.com/jiu-u/oai-api/pkg/adapter/api/v1"
	"io"
	"net/http"
)

type OpenAIProvider struct {
	Config
}

func NewOpenAIProvider(config Config) *OpenAIProvider {
	return &OpenAIProvider{
		Config: config,
	}
}

func (p *OpenAIProvider) ImageVariations(ctx context.Context, req *v1.CreateImageVariationRequest) (io.Reader, http.Header, error) {
	url := p.Config.EndPoint + "/v1/images/variations"
	bodyBytes, err := sonic.Marshal(req)
	if err != nil {
		return nil, nil, err
	}
	return p.DoFormRequest(ctx, url, bytes.NewBuffer(bodyBytes))
}

func (p *OpenAIProvider) ImageVariationsByBytes(ctx context.Context, req []byte) (io.Reader, http.Header, error) {
	url := p.Config.EndPoint + "/v1/images/variations"
	return p.DoFormRequest(ctx, url, bytes.NewBuffer(req))
}

func (p *OpenAIProvider) CreateImageEdit(ctx context.Context, req *v1.EditImageRequest) (io.Reader, http.Header, error) {
	url := p.Config.EndPoint + "/v1/images/edits"
	bodyBytes, err := sonic.Marshal(req)
	if err != nil {
		return nil, nil, err
	}
	return p.DoFormRequest(ctx, url, bytes.NewBuffer(bodyBytes))
}

func (p *OpenAIProvider) CreateImageEditByBytes(ctx context.Context, req []byte) (io.Reader, http.Header, error) {
	url := p.Config.EndPoint + "/v1/images/edits"
	return p.DoFormRequest(ctx, url, bytes.NewBuffer(req))
}

func (p *OpenAIProvider) ChatCompletions(ctx context.Context, req *v1.ChatCompletionRequest) (io.Reader, http.Header, error) {
	bodyBytes, err := sonic.Marshal(req)
	if err != nil {
		return nil, nil, err
	}
	return p.DoChatCompletionsRequest(ctx, bytes.NewBuffer(bodyBytes))
}

func (p *OpenAIProvider) ChatCompletionsByBytes(ctx context.Context, req []byte) (io.Reader, http.Header, error) {
	return p.DoChatCompletionsRequest(ctx, bytes.NewBuffer(req))
}

func (p *OpenAIProvider) Models(ctx context.Context) ([]string, error) {
	url := p.Config.EndPoint + "/v1/models"
	respBody, _, err := p.DoJsonRequest(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(respBody)
	if err != nil {
		return nil, err
	}
	var resp v1.ModelResp
	err = sonic.Unmarshal(bodyBytes, &resp)
	if err != nil {
		return nil, err
	}
	var models []string
	for _, model := range resp.Data {
		models = append(models, model.ID)
	}
	return models, nil
}

func (p *OpenAIProvider) Completions(ctx context.Context, req *v1.CompletionsRequest) (io.Reader, http.Header, error) {
	bodyBytes, err := sonic.Marshal(req)
	if err != nil {
		return nil, nil, err
	}
	return p.DoCompletionsRequest(ctx, bytes.NewBuffer(bodyBytes))
}

func (p *OpenAIProvider) CompletionsByBytes(ctx context.Context, req []byte) (io.Reader, http.Header, error) {
	return p.DoCompletionsRequest(ctx, bytes.NewBuffer(req))
}

func (p *OpenAIProvider) Embeddings(ctx context.Context, req *v1.EmbeddingRequest) (io.Reader, http.Header, error) {
	bodyBytes, err := sonic.Marshal(req)
	if err != nil {
		return nil, nil, err
	}
	return p.DoEmbeddingsRequest(ctx, bytes.NewBuffer(bodyBytes))
}

func (p *OpenAIProvider) EmbeddingsByBytes(ctx context.Context, req []byte) (io.Reader, http.Header, error) {
	return p.DoEmbeddingsRequest(ctx, bytes.NewBuffer(req))
}

func (p *OpenAIProvider) CreateSpeech(ctx context.Context, req *v1.SpeechRequest) (io.Reader, http.Header, error) {
	bodyBytes, err := sonic.Marshal(req)
	if err != nil {
		return nil, nil, err
	}
	return p.DoTextToSpeechRequest(ctx, bytes.NewBuffer(bodyBytes))
}

func (p *OpenAIProvider) CreateSpeechByBytes(ctx context.Context, req []byte) (io.Reader, http.Header, error) {
	return p.DoTextToSpeechRequest(ctx, bytes.NewBuffer(req))
}

func (p *OpenAIProvider) Transcriptions(ctx context.Context, req *v1.TranscriptionRequest) (io.Reader, http.Header, error) {
	url := p.Config.EndPoint + "/v1/audio/transcriptions"
	bodyBytes, err := sonic.Marshal(req)
	if err != nil {
		return nil, nil, err
	}
	return p.DoFormRequest(ctx, url, bytes.NewBuffer(bodyBytes))
}

func (p *OpenAIProvider) TranscriptionsByBytes(ctx context.Context, req []byte) (io.Reader, http.Header, error) {
	url := p.Config.EndPoint + "/v1/audio/transcriptions"
	return p.DoFormRequest(ctx, url, bytes.NewBuffer(req))
}

func (p *OpenAIProvider) Translations(ctx context.Context, req *v1.TranslationRequest) (io.Reader, http.Header, error) {
	url := p.Config.EndPoint + "/v1/audio/translations"
	bodyBytes, err := sonic.Marshal(req)
	if err != nil {
		return nil, nil, err
	}
	return p.DoFormRequest(ctx, url, bytes.NewBuffer(bodyBytes))
}

func (p *OpenAIProvider) TranslationsByBytes(ctx context.Context, req []byte) (io.Reader, http.Header, error) {
	url := p.Config.EndPoint + "/v1/audio/translations"
	return p.DoFormRequest(ctx, url, bytes.NewBuffer(req))
}

func (p *OpenAIProvider) CreateImage(ctx context.Context, req *v1.CreateImageRequest) (io.Reader, http.Header, error) {
	url := p.Config.EndPoint + "/v1/images/generations"
	bodyBytes, err := sonic.Marshal(req)
	if err != nil {
		return nil, nil, err
	}
	return p.DoJsonRequest(ctx, url, bytes.NewBuffer(bodyBytes))
}

func (p *OpenAIProvider) CreateImageByBytes(ctx context.Context, req []byte) (io.Reader, http.Header, error) {
	url := p.Config.EndPoint + "/v1/images/generations"
	return p.DoJsonRequest(ctx, url, bytes.NewBuffer(req))
}

// -------------

func (p *OpenAIProvider) DoJsonRequest(ctx context.Context, url string, body io.Reader) (io.Reader, http.Header, error) {
	request, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+p.Config.APIKey)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	return resp.Body, resp.Header, nil
}

func (p *OpenAIProvider) DoRequest(ctx context.Context, url string, body io.Reader, contextType string) (io.Reader, http.Header, error) {
	request, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, nil, err
	}
	request.Header.Set("Content-Type", contextType)
	request.Header.Set("Authorization", "Bearer "+p.Config.APIKey)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	return resp.Body, resp.Header, nil
}

func (p *OpenAIProvider) DoFormRequest(ctx context.Context, url string, body io.Reader) (io.Reader, http.Header, error) {
	// 创建 multipart 写入器
	return p.DoRequest(ctx, url, body, "application/x-www-form-urlencoded")
}

func (p *OpenAIProvider) DoChatCompletionsRequest(ctx context.Context, input io.Reader) (io.Reader, http.Header, error) {
	url := p.Config.EndPoint + "/v1/chat/completions"
	return p.DoJsonRequest(ctx, url, input)
}

func (p *OpenAIProvider) DoEmbeddingsRequest(ctx context.Context, input io.Reader) (io.Reader, http.Header, error) {
	url := p.Config.EndPoint + "/v1/embeddings"
	return p.DoJsonRequest(ctx, url, input)
}

func (p *OpenAIProvider) DoCompletionsRequest(ctx context.Context, input io.Reader) (io.Reader, http.Header, error) {
	url := p.Config.EndPoint + "/v1/completions"
	return p.DoJsonRequest(ctx, url, input)
}

func (p *OpenAIProvider) DoTextToSpeechRequest(ctx context.Context, input io.Reader) (io.Reader, http.Header, error) {
	url := p.Config.EndPoint + "/v1/audio/speech"
	return p.DoJsonRequest(ctx, url, input)
}
