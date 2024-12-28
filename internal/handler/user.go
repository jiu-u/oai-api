package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jiu-u/oai-api/internal/service"
)

type UserHandler struct {
	*Handler
	svc service.UserService
}

func NewUserHandler(handler *Handler, svc service.UserService) *UserHandler {
	return &UserHandler{
		Handler: handler,
		svc:     svc,
	}
}

func (h *UserHandler) GetUser(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	resp, err := h.svc.GetUserInfo(ctx, userId)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, resp)
}
