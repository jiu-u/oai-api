package handler

import (
	"errors"
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
	if userId == 0 {
		apiV1.HandleError(ctx, 500, errors.New("userId is required"), "userId is required")
		return
	}
	req := new(apiV1.ResetApiKeyRequest)
	req.UserId = strconv.FormatUint(userId, 10)
	resp, err := h.svc.ResetApiKey(ctx, req)
	if err != nil {
		apiV1.HandleError(ctx, 400, err, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, resp)
}

func (h *ApiKeyHandler) GetApiKey(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == 0 {
		apiV1.HandleError(ctx, 500, errors.New("userId is required"), "userId is required")
		return
	}
	resp, err := h.svc.GetUserApiKey(ctx, userId)
	if err != nil {
		apiV1.HandleError(ctx, 400, err, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, resp)
}
