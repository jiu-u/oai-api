package handler

import (
	"github.com/gin-gonic/gin"
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/service"
)

type VerificationHandler struct {
	*Handler
	svc service.VerificationService
}

func NewVerificationHandler(handler *Handler, svc service.VerificationService) *VerificationHandler {
	return &VerificationHandler{
		Handler: handler,
		svc:     svc,
	}
}

func (h *VerificationHandler) SetVerificationCode2Email(ctx *gin.Context) {
	req := new(apiV1.VerificationEmailReq)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if req.Email != req.Code {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, "email and code not match")
		return
	}
	err := h.svc.SendEmailVerificationCode(ctx, req.Email)
	if err != nil {
		apiV1.HandleError(ctx, 0, err, nil)
		return
	}
	apiV1.HandleSuccess(ctx, nil)
}
