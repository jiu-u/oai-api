package provider

import (
	"bytes"
	"context"
	"errors"
	"github.com/bytedance/sonic"
	v1 "github.com/jiu-u/oai-api/pkg/adapter/api/v1"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type OpenAIProvider struct {
	Config
}

func NewOpenAIProvider(config Config) *OpenAIProvider {
	return &OpenAIProvider{
		Config: config,
	}
}

func (p *OpenAIProvider) ChatCompletions(ctx context.Context, req *v1.ChatCompletionRequest) (io.ReadCloser, http.Header, error) {
	bodyBytes, err := sonic.Marshal(req)
	if err != nil {
		return nil, nil, err
	}
	return p.DoChatCompletionsRequest(ctx, bytes.NewBuffer(bodyBytes))
}

func (p *OpenAIProvider) ChatCompletionsByBytes(ctx context.Context, req []byte) (io.ReadCloser, http.Header, error) {
	return p.DoChatCompletionsRequest(ctx, bytes.NewBuffer(req))
}

func (p *OpenAIProvider) Models(ctx context.Context) ([]string, error) {
	url := p.Config.EndPoint + "/v1/models"
	respBody, _, err := p.DoJsonRequest(ctx, "GET", url, nil)
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

func (p *OpenAIProvider) Completions(ctx context.Context, req *v1.CompletionsRequest) (io.ReadCloser, http.Header, error) {
	bodyBytes, err := sonic.Marshal(req)
	if err != nil {
		return nil, nil, err
	}
	return p.DoCompletionsRequest(ctx, bytes.NewBuffer(bodyBytes))
}

func (p *OpenAIProvider) CompletionsByBytes(ctx context.Context, req []byte) (io.ReadCloser, http.Header, error) {
	return p.DoCompletionsRequest(ctx, bytes.NewBuffer(req))
}

func (p *OpenAIProvider) Embeddings(ctx context.Context, req *v1.EmbeddingRequest) (io.ReadCloser, http.Header, error) {
	bodyBytes, err := sonic.Marshal(req)
	if err != nil {
		return nil, nil, err
	}
	return p.DoEmbeddingsRequest(ctx, bytes.NewBuffer(bodyBytes))
}

func (p *OpenAIProvider) EmbeddingsByBytes(ctx context.Context, req []byte) (io.ReadCloser, http.Header, error) {
	return p.DoEmbeddingsRequest(ctx, bytes.NewBuffer(req))
}

func (p *OpenAIProvider) CreateSpeech(ctx context.Context, req *v1.SpeechRequest) (io.ReadCloser, http.Header, error) {
	bodyBytes, err := sonic.Marshal(req)
	if err != nil {
		return nil, nil, err
	}
	return p.DoTextToSpeechRequest(ctx, bytes.NewBuffer(bodyBytes))
}

func (p *OpenAIProvider) CreateSpeechByBytes(ctx context.Context, req []byte) (io.ReadCloser, http.Header, error) {
	return p.DoTextToSpeechRequest(ctx, bytes.NewBuffer(req))
}

func (p *OpenAIProvider) Transcriptions(ctx context.Context, req *v1.TranscriptionRequest) (io.ReadCloser, http.Header, error) {
	url := p.Config.EndPoint + "/v1/audio/transcriptions"
	// 创建一个字节缓冲区来存储请求体
	var buf bytes.Buffer
	// 创建 multipart 写入器
	writer := multipart.NewWriter(&buf)
	// 添加文件字段
	part, err := writer.CreateFormFile("file", req.File.Filename)
	if err != nil {
		return nil, nil, err
	}

	file, err := req.File.Open()
	if err != nil {
		return nil, nil, err
	}
	defer file.Close() // 确保在函数结束时关闭文件

	// 将文件内容复制到文件字段
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, nil, err
	}

	// 添加其他字段并检查错误
	if err := writer.WriteField("model", req.Model); err != nil {
		return nil, nil, err
	}
	if req.Prompt != "" {
		if err := writer.WriteField("prompt", req.Prompt); err != nil {
			return nil, nil, err
		}
	}
	if req.ResponseFormat != "" {
		if err := writer.WriteField("response_format", req.ResponseFormat); err != nil {
			return nil, nil, err
		}
	}
	if req.Temperature != 0 {
		if err := writer.WriteField("temperature", strconv.FormatFloat(req.Temperature, 'f', -1, 64)); err != nil {
			return nil, nil, err
		}
	}
	// 关闭 multipart 写入器
	err = writer.Close()
	if err != nil {
		return nil, nil, err
	}
	// 返回请求
	return p.DoFormRequest(ctx, url, &buf, writer.FormDataContentType())
}

func (p *OpenAIProvider) Translations(ctx context.Context, req *v1.TranslationRequest) (io.ReadCloser, http.Header, error) {
	url := p.Config.EndPoint + "/v1/audio/translations"
	// 创建一个字节缓冲区来存储请求体
	var buf bytes.Buffer
	// 创建 multipart 写入器
	writer := multipart.NewWriter(&buf)
	// 添加文件字段
	part, err := writer.CreateFormFile("file", req.File.Filename)
	if err != nil {
		return nil, nil, err
	}

	file, err := req.File.Open()
	if err != nil {
		return nil, nil, err
	}
	defer file.Close() // 确保在函数结束时关闭文件

	// 将文件内容复制到文件字段
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, nil, err
	}

	// 添加其他字段并检查错误
	if err := writer.WriteField("model", req.Model); err != nil {
		return nil, nil, err
	}
	if req.Prompt != "" {
		if err := writer.WriteField("prompt", req.Prompt); err != nil {
			return nil, nil, err
		}
	}
	if req.ResponseFormat != "" {
		if err := writer.WriteField("response_format", req.ResponseFormat); err != nil {
			return nil, nil, err
		}
	}
	if req.Temperature != 0 {
		if err := writer.WriteField("temperature", strconv.FormatFloat(req.Temperature, 'f', -1, 64)); err != nil {
			return nil, nil, err
		}
	}
	// 关闭 multipart 写入器
	err = writer.Close()
	if err != nil {
		return nil, nil, err
	}
	// 返回请求
	return p.DoFormRequest(ctx, url, &buf, writer.FormDataContentType())
}

func (p *OpenAIProvider) CreateImage(ctx context.Context, req *v1.CreateImageRequest) (io.ReadCloser, http.Header, error) {
	url := p.Config.EndPoint + "/v1/images/generations"
	bodyBytes, err := sonic.Marshal(req)
	if err != nil {
		return nil, nil, err
	}
	return p.DoJsonRequest(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
}

func (p *OpenAIProvider) CreateImageByBytes(ctx context.Context, req []byte) (io.ReadCloser, http.Header, error) {
	url := p.Config.EndPoint + "/v1/images/generations"
	return p.DoJsonRequest(ctx, "POST", url, bytes.NewBuffer(req))
}

func (p *OpenAIProvider) CreateImageEdit(ctx context.Context, req *v1.EditImageRequest) (io.ReadCloser, http.Header, error) {
	url := p.Config.EndPoint + "/v1/images/edits"

	// 创建一个字节缓冲区来存储请求体
	var buf bytes.Buffer
	// 创建 multipart 写入器
	writer := multipart.NewWriter(&buf)
	// 添加文件字段
	part, err := writer.CreateFormFile("file", req.Image.Filename)
	if err != nil {
		return nil, nil, err
	}

	file, err := req.Image.Open()
	if err != nil {
		return nil, nil, err
	}
	defer file.Close() // 确保在函数结束时关闭文件
	// 将文件内容复制到文件字段
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, nil, err
	}

	if req.Mask != nil && req.Mask.Filename != "" && req.Mask.Size > 0 {
		maskPart, err := writer.CreateFormFile("mask", req.Mask.Filename)
		if err != nil {
			return nil, nil, err
		}

		maskFile, err := req.Mask.Open()
		if err != nil {
			return nil, nil, err
		}
		defer maskFile.Close() // 确保在函数结束时关闭文件
		_, err = io.Copy(maskPart, maskFile)
		if err != nil {
			return nil, nil, err
		}
	}

	// 添加其他字段并检查错误
	if req.Model == "" {
		if err := writer.WriteField("model", req.Model); err != nil {
			return nil, nil, err
		}
	}
	if req.Prompt != "" {
		if err := writer.WriteField("prompt", req.Prompt); err != nil {
			return nil, nil, err
		}
	}
	if req.ResponseFormat != "" {
		if err := writer.WriteField("response_format", req.ResponseFormat); err != nil {
			return nil, nil, err
		}
	}
	if req.N != 0 {
		if err := writer.WriteField("temperature", strconv.FormatInt(int64(req.N), 10)); err != nil {
			return nil, nil, err
		}
	}
	if req.Size != "" {
		if err := writer.WriteField("size", req.Size); err != nil {
			return nil, nil, err
		}
	}
	if req.User != "" {
		if err := writer.WriteField("user", req.User); err != nil {
			return nil, nil, err
		}
	}
	// 关闭 multipart 写入器
	err = writer.Close()
	if err != nil {
		return nil, nil, err
	}
	// 返回请求
	return p.DoFormRequest(ctx, url, &buf, writer.FormDataContentType())
}

func (p *OpenAIProvider) ImageVariations(ctx context.Context, req *v1.CreateImageVariationRequest) (io.ReadCloser, http.Header, error) {
	url := p.Config.EndPoint + "/v1/images/variations"
	// 创建一个字节缓冲区来存储请求体
	var buf bytes.Buffer
	// 创建 multipart 写入器
	writer := multipart.NewWriter(&buf)
	// 添加文件字段
	part, err := writer.CreateFormFile("file", req.Image.Filename)
	if err != nil {
		return nil, nil, err
	}

	file, err := req.Image.Open()
	if err != nil {
		return nil, nil, err
	}
	defer file.Close() // 确保在函数结束时关闭文件

	// 将文件内容复制到文件字段
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, nil, err
	}

	// 添加其他字段并检查错误
	if req.Model == "" {
		if err := writer.WriteField("model", req.Model); err != nil {
			return nil, nil, err
		}
	}

	if req.ResponseFormat != "" {
		if err := writer.WriteField("response_format", req.ResponseFormat); err != nil {
			return nil, nil, err
		}
	}
	if req.N != 0 {
		if err := writer.WriteField("temperature", strconv.FormatInt(int64(req.N), 10)); err != nil {
			return nil, nil, err
		}
	}
	if req.Size != "" {
		if err := writer.WriteField("size", req.Size); err != nil {
			return nil, nil, err
		}
	}
	if req.User != "" {
		if err := writer.WriteField("user", req.User); err != nil {
			return nil, nil, err
		}
	}

	// 关闭 multipart 写入器
	err = writer.Close()
	if err != nil {
		return nil, nil, err
	}
	// 返回请求
	return p.DoFormRequest(ctx, url, &buf, writer.FormDataContentType())
}

// -------------

func (p *OpenAIProvider) DoJsonRequest(ctx context.Context, Method string, url string, body io.Reader) (io.ReadCloser, http.Header, error) {
	return p.DoRequest(ctx, url, Method, body, "application/json")
	//request, err := http.NewRequestWithContext(ctx, "POST", url, body)
	//if err != nil {
	//	return nil, nil, err
	//}
	//request.Header.Set("Content-Type", "application/json")
	//request.Header.Set("Authorization", "Bearer "+p.Config.APIKey)
	//client := &http.Client{}
	//resp, err := client.Do(request)
	//if err != nil {
	//	return nil, nil, err
	//}
	//if resp.StatusCode != 200 {
	//	return nil, nil, errors.New(resp.Status)
	//}
	//return resp.Body, resp.Header, nil
}

func (p *OpenAIProvider) DoRequest(ctx context.Context, url string, Method string, body io.Reader, contextType string) (io.ReadCloser, http.Header, error) {
	request, err := http.NewRequestWithContext(ctx, Method, url, body)
	if err != nil {
		return nil, nil, err
	}
	request.Header.Set("Content-Type", contextType)
	request.Header.Set("Authorization", "Bearer "+p.Config.APIKey)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != 200 {
		return nil, nil, errors.New(resp.Status)
	}
	return resp.Body, resp.Header, nil
}

func (p *OpenAIProvider) DoFormRequest(ctx context.Context, url string, body io.Reader, contentType string) (io.ReadCloser, http.Header, error) {
	// 创建 multipart 写入器
	return p.DoRequest(ctx, url, "POST", body, contentType)
}

func (p *OpenAIProvider) DoChatCompletionsRequest(ctx context.Context, input io.Reader) (io.ReadCloser, http.Header, error) {
	url := p.Config.EndPoint + "/v1/chat/completions"
	return p.DoJsonRequest(ctx, "POST", url, input)
}

func (p *OpenAIProvider) DoEmbeddingsRequest(ctx context.Context, input io.Reader) (io.ReadCloser, http.Header, error) {
	url := p.Config.EndPoint + "/v1/embeddings"
	return p.DoJsonRequest(ctx, "POST", url, input)
}

func (p *OpenAIProvider) DoCompletionsRequest(ctx context.Context, input io.Reader) (io.ReadCloser, http.Header, error) {
	url := p.Config.EndPoint + "/v1/completions"
	return p.DoJsonRequest(ctx, "POST", url, input)
}

func (p *OpenAIProvider) DoTextToSpeechRequest(ctx context.Context, input io.Reader) (io.ReadCloser, http.Header, error) {
	url := p.Config.EndPoint + "/v1/audio/speech"
	return p.DoJsonRequest(ctx, "POST", url, input)
}
