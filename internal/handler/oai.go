package handler

import (
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
	for k, v := range respHeader {
		ctx.Header(k, v[0])
	}
	
	ctx.Data(http.StatusOK, contentType, respBytes)
}
