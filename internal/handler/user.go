package handler

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/jiu-u/oai-api/api/v1"
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

func (h *UserHandler) GetCurrentUser(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	resp, err := h.svc.GetUserInfo(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, 400, err, err.Error())
		return
	}
	v1.HandleSuccess(ctx, resp)
}
