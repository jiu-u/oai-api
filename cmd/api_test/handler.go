package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	v1 "github.com/jiu-u/oai-api/pkg/adapter/api/v1"
	"github.com/jiu-u/oai-api/pkg/adapter/provider"
	"io"
	"net/http"
	"time"
)

type OpenAIHandler struct {
	Provider provider.Provider
}

func NewOpenAIHandler(provider provider.Provider) *OpenAIHandler {
	return &OpenAIHandler{Provider: provider}
}

func HandleResponse1(ctx *gin.Context, responseBody io.ReadCloser, respHeader http.Header) {
	defer responseBody.Close()
	for k, v := range respHeader {
		ctx.Writer.Header().Set(k, v[0])
	}
	contentType := respHeader.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}
	ctx.Writer.Header().Set("Content-Type", contentType)
	_, err := io.Copy(ctx.Writer, responseBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
	}
}

func HandleResponse2(ctx *gin.Context, responseBody io.ReadCloser, respHeader http.Header) {
	defer responseBody.Close()
	for k, v := range respHeader {
		ctx.Writer.Header().Set(k, v[0])
	}
	contentType := respHeader.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}
	reader := bufio.NewReader(responseBody)
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			if _, writeErr := ctx.Writer.Write(buf[:n]); writeErr != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
				return
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
			return
		}
	}
}

func HandleResponse3(ctx *gin.Context, responseBody io.ReadCloser, respHeader http.Header) {
	defer responseBody.Close()
	for k, v := range respHeader {
		ctx.Writer.Header().Set(k, v[0])
	}
	contentType := respHeader.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}
	scanner := bufio.NewScanner(responseBody)
	for scanner.Scan() {
		line := scanner.Bytes()
		if _, err := ctx.Writer.Write(line); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
			return
		}
		ctx.Writer.Write([]byte("\n"))
		ctx.Writer.Flush()
	}
	if err := scanner.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}
}

func HandleResponse4(ctx *gin.Context, responseBody io.ReadCloser, respHeader http.Header) {
	defer responseBody.Close()
	for k, v := range respHeader {
		ctx.Writer.Header().Set(k, v[0])
	}
	contentType := respHeader.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}
	ctx.Stream(func(w io.Writer) bool {
		scanner := bufio.NewScanner(responseBody)
		for scanner.Scan() {
			line := scanner.Bytes()
			if _, err := w.Write(line); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
				return false
			}
			w.Write([]byte("\n"))
			w.(http.Flusher).Flush()
		}
		if err := scanner.Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
			return false
		}
		return false
	})

}

func HandleResponse5(ctx *gin.Context, responseBody io.ReadCloser, respHeader http.Header) {
	defer responseBody.Close()
	for k, v := range respHeader {
		ctx.Writer.Header().Set(k, v[0])
	}
	contentType := respHeader.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}
	if contentType != "text/event-stream" {
		HandleResponse1(ctx, responseBody, respHeader)
		return
	}
	scanner := bufio.NewScanner(responseBody)
	for scanner.Scan() {
		line := scanner.Bytes()
		res := make(map[string]any)
		if err := sonic.Unmarshal(line, &res); err != nil {
			fmt.Println("line", string(line))
		}
		jsonStr, _ := sonic.Marshal(res["data"])
		ctx.SSEvent("data", jsonStr)
		ctx.Writer.Flush()
	}
	if err := scanner.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}
}

func (h *OpenAIHandler) ChatCompletions(ctx *gin.Context) {
	start := time.Now()
	fmt.Println("time", start.Format("2006-01-02 15:04:05"))
	var req v1.ChatCompletionRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	responseBody, respHeader, err := h.Provider.ChatCompletions(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("time", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("duration", time.Since(start))
	fmt.Println(333)
	HandleResponse3(ctx, responseBody, respHeader)
	//// 读取响应数据
	//for k, v := range respHeader {
	//	ctx.Header(k, v[0])
	//}
	//respBytes, err := io.ReadAll(responseBody)
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
	//	return
	//}
	//
	//// 设置响应头并返回 OpenAI API 的响应
	////ctx.Header("Content-Type", respHeader.Get("Content-Type"))
	//for k, v := range respHeader {
	//	ctx.Writer.Header().Set(k, v[0])
	//}
	//ctx.Data(http.StatusOK, "application/json", respBytes)
}

func (h *OpenAIHandler) ChatCompletionsByBytes(ctx *gin.Context) {
	var req v1.ChatCompletionRequest
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//bodyBytes, err = changeBytesModelId(bodyBytes, "deepseek-chat")
	responseBody, respHeader, err := h.Provider.ChatCompletionsByBytes(ctx, bodyBytes)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 读取响应数据
	respBytes, err := io.ReadAll(responseBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}

	// 设置响应头并返回 OpenAI API 的响应
	//ctx.Header("Content-Type", respHeader.Get("Content-Type"))
	for k, v := range respHeader {
		ctx.Header(k, v[0])
	}
	ctx.Data(http.StatusOK, "application/json", respBytes)
}

func (h *OpenAIHandler) Completions(ctx *gin.Context) {
	var req v1.CompletionsRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	responseBody, respHeader, err := h.Provider.Completions(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 读取响应数据
	respBytes, err := io.ReadAll(responseBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}
	// 设置响应头并返回 OpenAI API 的响应
	//ctx.Header("Content-Type", respHeader.Get("Content-Type"))
	for k, v := range respHeader {
		ctx.Header(k, v[0])
	}
	ctx.Data(http.StatusOK, "application/json", respBytes)
}

func (h *OpenAIHandler) CompletionsByBytes(ctx *gin.Context) {
	var req v1.CompletionsRequest
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	responseBody, respHeader, err := h.Provider.CompletionsByBytes(ctx, bodyBytes)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 读取响应数据
	respBytes, err := io.ReadAll(responseBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}
	// 设置响应头并返回 OpenAI API 的响应
	//ctx.Header("Content-Type", respHeader.Get("Content-Type"))
	for k, v := range respHeader {
		ctx.Header(k, v[0])
	}
	ctx.Data(http.StatusOK, "application/json", respBytes)
}

func GetModelRespByModels(models []string) *v1.ModelResp {
	resp := &v1.ModelResp{
		Object: "list",
		Data:   make([]v1.Model, len(models)),
	}
	for idx, modelId := range models {
		resp.Data[idx] = v1.Model{
			ID:         modelId,
			Object:     "model",
			Created:    0,
			OwnedBy:    "system",
			Permission: nil,
		}
	}
	return resp
}

func (h *OpenAIHandler) Models(ctx *gin.Context) {
	models, err := h.Provider.Models(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	resp := GetModelRespByModels(models)
	ctx.JSON(200, resp)
}

func (h *OpenAIHandler) Embeddings(ctx *gin.Context) {
	var req v1.EmbeddingRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	responseBody, respHeader, err := h.Provider.Embeddings(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 读取响应数据
	respBytes, err := io.ReadAll(responseBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}
	// 设置响应头并返回 OpenAI API 的响应
	//ctx.Header("Content-Type", respHeader.Get("Content-Type"))
	for k, v := range respHeader {
		ctx.Header(k, v[0])
	}
	ctx.Data(http.StatusOK, "application/json", respBytes)
}

func (h *OpenAIHandler) EmbeddingsByBytes(ctx *gin.Context) {
	var req v1.EmbeddingRequest
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	responseBody, respHeader, err := h.Provider.EmbeddingsByBytes(ctx, bodyBytes)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 读取响应数据
	respBytes, err := io.ReadAll(responseBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}
	// 设置响应头并返回 OpenAI API 的响应
	//ctx.Header("Content-Type", respHeader.Get("Content-Type"))
	for k, v := range respHeader {
		ctx.Header(k, v[0])
	}
	ctx.Data(http.StatusOK, "application/json", respBytes)
}

func (h *OpenAIHandler) AudioSpeech(ctx *gin.Context) {
	var req v1.SpeechRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	responseBody, respHeader, err := h.Provider.CreateSpeech(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 读取响应数据
	respBytes, err := io.ReadAll(responseBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}
	// 设置响应头并返回 OpenAI API 的响应
	//ctx.Header("Content-Type", respHeader.Get("Content-Type"))
	for k, v := range respHeader {
		ctx.Header(k, v[0])
	}
	ctx.Data(http.StatusOK, "audio/mpeg", respBytes)
}

func (h *OpenAIHandler) AudioSpeechByBytes(ctx *gin.Context) {
	var req v1.SpeechRequest
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	responseBody, respHeader, err := h.Provider.CreateSpeechByBytes(ctx, bodyBytes)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 读取响应数据
	respBytes, err := io.ReadAll(responseBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}
	// 设置响应头并返回 OpenAI API 的响应
	//ctx.Header("Content-Type", respHeader.Get("Content-Type"))
	for k, v := range respHeader {
		ctx.Header(k, v[0])
	}
	ctx.Data(http.StatusOK, "audio/mpeg", respBytes)
}

func (h *OpenAIHandler) AudioTranscription(ctx *gin.Context) {
	var req v1.TranscriptionRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("req: %+v\n", req)
	responseBody, respHeader, err := h.Provider.Transcriptions(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 读取响应数据
	respBytes, err := io.ReadAll(responseBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}
	// 设置响应头并返回 OpenAI API 的响应
	//ctx.Header("Content-Type", respHeader.Get("Content-Type"))
	for k, v := range respHeader {
		ctx.Header(k, v[0])
	}
	ctx.Data(http.StatusOK, "application/json", respBytes)
}

func (h *OpenAIHandler) AudioTranslation(ctx *gin.Context) {
	var req v1.TranslationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	responseBody, respHeader, err := h.Provider.Translations(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 读取响应数据
	respBytes, err := io.ReadAll(responseBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}
	// 设置响应头并返回 OpenAI API 的响应
	contentType := respHeader.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}
	//ctx.Header("Content-Type", respHeader.Get("Content-Type"))
	for k, v := range respHeader {
		ctx.Header(k, v[0])
	}
	ctx.Data(http.StatusOK, contentType, respBytes)
}

func (h *OpenAIHandler) ImageGeneration(ctx *gin.Context) {
	var req v1.CreateImageRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	responseBody, respHeader, err := h.Provider.CreateImage(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 读取响应数据
	respBytes, err := io.ReadAll(responseBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}
	// 设置响应头并返回 OpenAI API 的响应
	contentType := respHeader.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}
	//ctx.Header("Content-Type", respHeader.Get("Content-Type"))
	for k, v := range respHeader {
		ctx.Header(k, v[0])
	}
	ctx.Data(http.StatusOK, contentType, respBytes)
}

func (h *OpenAIHandler) ImageGenerationByBytes(ctx *gin.Context) {
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	var req v1.CreateImageRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	responseBody, respHeader, err := h.Provider.CreateImageByBytes(ctx, bodyBytes)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 读取响应数据
	respBytes, err := io.ReadAll(responseBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}
	// 设置响应头并返回 OpenAI API 的响应
	contentType := respHeader.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}
	//ctx.Header("Content-Type", respHeader.Get("Content-Type"))
	for k, v := range respHeader {
		ctx.Header(k, v[0])
	}
	ctx.Data(http.StatusOK, contentType, respBytes)
}

func (h *OpenAIHandler) ImageEdit(ctx *gin.Context) {
	var req v1.EditImageRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	responseBody, respHeader, err := h.Provider.CreateImageEdit(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 读取响应数据
	respBytes, err := io.ReadAll(responseBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}
	// 设置响应头并返回 OpenAI API 的响应
	contentType := respHeader.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}
	//ctx.Header("Content-Type", respHeader.Get("Content-Type"))
	for k, v := range respHeader {
		ctx.Header(k, v[0])
	}
	ctx.Data(http.StatusOK, contentType, respBytes)
}

func (h *OpenAIHandler) ImageVariation(ctx *gin.Context) {
	var req v1.CreateImageVariationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	responseBody, respHeader, err := h.Provider.ImageVariations(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 读取响应数据
	respBytes, err := io.ReadAll(responseBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}
	// 设置响应头并返回 OpenAI API 的响应
	contentType := respHeader.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}
	//ctx.Header("Content-Type", respHeader.Get("Content-Type"))
	for k, v := range respHeader {
		ctx.Header(k, v[0])
	}
	ctx.Data(http.StatusOK, contentType, respBytes)
}

func ChangeModelId(req *v1.ChatCompletionRequest, newModelId string) {
	req.Model = newModelId
}

func changeBytesModelId(bodyBytes []byte, newModelId string) ([]byte, error) {
	var result map[string]any
	err := sonic.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}
	result["model"] = newModelId
	bytes, err := sonic.Marshal(result)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
