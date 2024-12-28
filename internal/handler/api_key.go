package handler

import (
	"github.com/gin-gonic/gin"
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/service"
	"strconv"
)

type ApiKeyHandler struct {
	Handler *Handler
	svc     service.ApiKeyService
}

func NewApiKeyHandler(handler *Handler, svc service.ApiKeyService) *ApiKeyHandler {
	return &ApiKeyHandler{
		Handler: handler,
		svc:     svc,
	}
}

func (h *ApiKeyHandler) ResetApiKey(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	req := new(apiV1.ResetApiKeyRequest)
	req.UserId = strconv.FormatUint(userId, 10)
	resp, err := h.svc.ResetApiKey(ctx, req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, resp)
}

func (h *ApiKeyHandler) GetApiKey(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	resp, err := h.svc.GetUserApiKey(ctx, userId)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, resp)
}
