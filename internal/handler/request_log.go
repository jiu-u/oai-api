package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/service"
)

type RequestLogHandler struct {
	Handler *Handler
	svc     service.RequestLogService
	limit   int
}

func NewRequestLogHandler(handler *Handler, svc service.RequestLogService) *RequestLogHandler {
	return &RequestLogHandler{
		Handler: handler,
		svc:     svc,
		limit:   30,
	}
}

func (h *RequestLogHandler) GetRequestLogs(ctx *gin.Context) {
	req := new(apiV1.RequestLogsQuery)
	if err := ctx.ShouldBind(req); err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}

	resp, err := h.svc.GetRequestLogs(ctx, req)
	if err != nil {
		apiV1.HandleError(ctx, 400, err, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, resp)
}

func (h *RequestLogHandler) GetUserRequestLogs(ctx *gin.Context) {
	req := new(apiV1.RequestLogsQuery)
	if err := ctx.ShouldBind(req); err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	userId := ctx.Param("userId")
	if userId == "" {
		apiV1.HandleError(ctx, 400, errors.New("userId is required"), "userId is required")
		return
	}
	req.UserId = userId
	resp, err := h.svc.GetRequestLogs(ctx, req)
	if err != nil {
		apiV1.HandleError(ctx, 400, err, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, resp)
}

func (h *RequestLogHandler) GetRequestLogsModelRanking(ctx *gin.Context) {
	req := new(apiV1.RequestLogsRankingRequest)
	if err := ctx.ShouldBind(req); err != nil {
		apiV1.HandleError(ctx, 400, err, err.Error())
		return
	}
	resp, err := h.svc.GetRequestLogsModelRanking(ctx, req)
	if err != nil {
		apiV1.HandleError(ctx, 400, err, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, resp)
}

func (h *RequestLogHandler) GetRequestLogsUserRanking(ctx *gin.Context) {
	req := new(apiV1.RequestLogsRankingRequest)
	if err := ctx.ShouldBind(req); err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	resp, err := h.svc.GetRequestLogsUserRanking(ctx, req)
	if err != nil {
		apiV1.HandleError(ctx, 400, err, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, resp)
}
