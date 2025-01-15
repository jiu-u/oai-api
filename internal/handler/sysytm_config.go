package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/dto"
	"github.com/jiu-u/oai-api/internal/service"
)

type SystemConfigHandler struct {
	*Handler
	svc service.SystemConfigService
}

func NewSystemConfigHandler(handler *Handler, svc service.SystemConfigService) *SystemConfigHandler {
	err := svc.InitSystemConfig(context.Background())
	if err != nil {
		panic(err)
	}
	return &SystemConfigHandler{
		Handler: handler,
		svc:     svc,
	}
}

func (h *SystemConfigHandler) SetEmailConfig(ctx *gin.Context) {
	req := new(apiV1.EmailConfig)
	if err := ctx.ShouldBind(req); err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}

	err := h.svc.SetEmailConfig(ctx, req)
	if err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, nil)
}

func (h *SystemConfigHandler) GetEmailConfig(ctx *gin.Context) {
	resp, err := h.svc.GetEmailConfig(ctx)
	if err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, resp)
}

func (h *SystemConfigHandler) IsEmailServiceAvailable(ctx *gin.Context) {
	resp, err := h.svc.IsEmailServiceAvailable(ctx)
	if err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, resp)
}

func (h *SystemConfigHandler) SetLinuxDoOAuthConfig(ctx *gin.Context) {
	req := new(apiV1.LinuxDoOAuthConfig)
	if err := ctx.ShouldBind(req); err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}

	err := h.svc.SetLinuxDoOAuthConfig(ctx, req)
	if err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, nil)
}

func (h *SystemConfigHandler) GetLinuxDoOAuthConfig(ctx *gin.Context) {
	resp, err := h.svc.GetLinuxDoOAuthConfig(ctx)
	if err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, resp)
}

func (h *SystemConfigHandler) IsLinuxDoOAuthServiceAvailable(ctx *gin.Context) {
	resp, err := h.svc.IsLinuxDoOAuthAvailable(ctx)
	if err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, resp)
}

func (h *SystemConfigHandler) SetGithubOAuthConfig(ctx *gin.Context) {
	req := new(apiV1.GithubOAuthConfig)
	if err := ctx.ShouldBind(req); err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}

	err := h.svc.SetGithubOAuthConfig(ctx, req)
	if err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, nil)
}

func (h *SystemConfigHandler) GetGithubOAuthConfig(ctx *gin.Context) {
	resp, err := h.svc.GetGithubOAuthConfig(ctx)
	if err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, resp)
}

func (h *SystemConfigHandler) IsGithubOAuthServiceAvailable(ctx *gin.Context) {
	resp, err := h.svc.IsGithubOAuthAvailable(ctx)
	if err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, resp)
}

func (h *SystemConfigHandler) SetRegisterConfig(c *gin.Context) {
	req := new(apiV1.RegisterConfig)
	if err := c.ShouldBind(req); err != nil {
		apiV1.HandleError(c, 400, apiV1.ErrBadRequest, err.Error())
		return
	}

	err := h.svc.SetRegisterConfig(c, req)
	if err != nil {
		apiV1.HandleError(c, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	apiV1.HandleSuccess(c, nil)
}

func (h *SystemConfigHandler) GetRegisterConfig(c *gin.Context) {
	resp, err := h.svc.GetRegisterConfig(c)
	if err != nil {
		apiV1.HandleError(c, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	apiV1.HandleSuccess(c, resp)
}

func (h *SystemConfigHandler) SetModelConfig(c *gin.Context) {
	var req dto.ModelConfig
	if err := c.ShouldBind(&req); err != nil {
		apiV1.HandleError(c, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	err := h.svc.SetModelConfig(c, &req)
	if err != nil {
		apiV1.HandleError(c, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	apiV1.HandleSuccess(c, nil)
}

func (h *SystemConfigHandler) GetModelConfig(c *gin.Context) {
	resp, err := h.svc.GetModelConfig(c)
	if err != nil {
		apiV1.HandleError(c, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	apiV1.HandleSuccess(c, resp)
}
