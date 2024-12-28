package handler

import (
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
	resp, err := h.svc.GetRealTimeData(ctx, h.limit)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, resp)
}

func (h *RequestLogHandler) GetRequestLogRanking(ctx *gin.Context) {
	req := new(apiV1.RequestLogRanking)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.svc.GetStatisticsData(ctx, req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, resp)
}
