package main

import (
	"bytes"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	v1 "github.com/jiu-u/oai-api/pkg/adapter/api/v1"
	"github.com/jiu-u/oai-api/pkg/adapter/provider"
	"io"
	"net/http"
)

type OpenAIHandler struct {
	Provider provider.Provider
}

func NewOpenAIHandler(provider provider.Provider) *OpenAIHandler {
	return &OpenAIHandler{Provider: provider}
}

func (h *OpenAIHandler) ChatCompletions(ctx *gin.Context) {
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
	// 读取响应数据
	respBytes, err := io.ReadAll(responseBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}

	// 设置响应头并返回 OpenAI API 的响应
	ctx.Header("Content-Type", respHeader.Get("Content-Type"))
	ctx.Data(http.StatusOK, "application/json", respBytes)
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
	ctx.Header("Content-Type", respHeader.Get("Content-Type"))
	ctx.Data(http.StatusOK, "application/json", respBytes)
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
