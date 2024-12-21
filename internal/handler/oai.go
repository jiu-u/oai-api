package handler

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jiu-u/oai-api/internal/service"
	v1 "github.com/jiu-u/oai-api/pkg/adapter/api/v1"
	"io"
	"net/http"
)

type OAIHandler struct {
	*Handler
	oaiService service.OaiService
}

func NewOAIHandler(oaiService service.OaiService) *OAIHandler {
	return &OAIHandler{
		oaiService: oaiService,
	}
}

var defaultRespHandle = HandleOAIResponse2

func (h *OAIHandler) ChatCompletions(ctx *gin.Context) {
	var req v1.ChatCompletionRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	responseBody, respHeader, err := h.oaiService.ChatCompletions(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 读取响应数据
	defaultRespHandle(ctx, responseBody, respHeader)
	//respBytes, err := io.ReadAll(responseBody)
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
	//	return
	//}
	//// 设置响应头并返回 OpenAI API 的响应
	//for k, v := range respHeader {
	//	ctx.Header(k, v[0])
	//}
	//
	//ctx.Data(http.StatusOK, getContentType(respHeader), respBytes)
}

func (h *OAIHandler) ChatCompletionsByBytes(ctx *gin.Context) {
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
	responseBody, respHeader, err := h.oaiService.ChatCompletionsByBytes(ctx, bodyBytes, req.Model)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defaultRespHandle(ctx, responseBody, respHeader)
}

func (h *OAIHandler) Completions(ctx *gin.Context) {
	var req v1.CompletionsRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	responseBody, respHeader, err := h.oaiService.Completions(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defaultRespHandle(ctx, responseBody, respHeader)
}

func (h *OAIHandler) CompletionsByBytes(ctx *gin.Context) {
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
	responseBody, respHeader, err := h.oaiService.CompletionsByBytes(ctx, bodyBytes, req.Model)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defaultRespHandle(ctx, responseBody, respHeader)
}

func (h *OAIHandler) Models(ctx *gin.Context) {
	resp, err := h.oaiService.Models(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 读取响应数据
	ctx.JSON(http.StatusOK, resp)
}

func (h *OAIHandler) Embeddings(ctx *gin.Context) {
	var req v1.EmbeddingRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	responseBody, respHeader, err := h.oaiService.Embeddings(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defaultRespHandle(ctx, responseBody, respHeader)
}

func (h *OAIHandler) EmbeddingsByBytes(ctx *gin.Context) {
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
	responseBody, respHeader, err := h.oaiService.EmbeddingsByBytes(ctx, bodyBytes, req.Model)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defaultRespHandle(ctx, responseBody, respHeader)
}

func (h *OAIHandler) AudioSpeech(ctx *gin.Context) {
	var req v1.SpeechRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	responseBody, respHeader, err := h.oaiService.CreateSpeech(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defaultRespHandle(ctx, responseBody, respHeader)
}

func (h *OAIHandler) AudioSpeechByBytes(ctx *gin.Context) {
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
	responseBody, respHeader, err := h.oaiService.CreateSpeechByBytes(ctx, bodyBytes, req.Model)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defaultRespHandle(ctx, responseBody, respHeader)
}

func (h *OAIHandler) AudioTranscription(ctx *gin.Context) {
	var req v1.TranscriptionRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("req: %+v\n", req)
	responseBody, respHeader, err := h.oaiService.Transcriptions(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defaultRespHandle(ctx, responseBody, respHeader)
}

func (h *OAIHandler) AudioTranslation(ctx *gin.Context) {
	var req v1.TranslationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	responseBody, respHeader, err := h.oaiService.Translations(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defaultRespHandle(ctx, responseBody, respHeader)
}

func (h *OAIHandler) ImageGeneration(ctx *gin.Context) {
	var req v1.CreateImageRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	responseBody, respHeader, err := h.oaiService.CreateImage(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defaultRespHandle(ctx, responseBody, respHeader)
}

func (h *OAIHandler) ImageGenerationByBytes(ctx *gin.Context) {
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

	responseBody, respHeader, err := h.oaiService.CreateImageByBytes(ctx, bodyBytes, req.Model)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defaultRespHandle(ctx, responseBody, respHeader)
}

func (h *OAIHandler) ImageEdit(ctx *gin.Context) {
	var req v1.EditImageRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	responseBody, respHeader, err := h.oaiService.CreateImageEdit(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defaultRespHandle(ctx, responseBody, respHeader)
}

func (h *OAIHandler) ImageVariation(ctx *gin.Context) {
	var req v1.CreateImageVariationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	responseBody, respHeader, err := h.oaiService.ImageVariations(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defaultRespHandle(ctx, responseBody, respHeader)
}
