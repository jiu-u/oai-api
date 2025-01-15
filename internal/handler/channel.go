package handler

import (
	"github.com/gin-gonic/gin"
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/service"
	"strconv"
)

type ChannelHandler struct {
	*Handler
	svc      service.ChannelService
	checkSvc service.ModelCheckService
}

func NewChannelHandler(
	handler *Handler,
	svc service.ChannelService,
	checkSvc service.ModelCheckService,
) *ChannelHandler {
	return &ChannelHandler{
		Handler:  handler,
		svc:      svc,
		checkSvc: checkSvc,
	}
}

func (h *ChannelHandler) GetChannels(ctx *gin.Context) {
	req := new(apiV1.ChannelQueryRequest)
	if err := ctx.ShouldBind(req); err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	resp, err := h.svc.GetChannels(ctx, req)
	if err != nil {
		apiV1.HandleError(ctx, 0, err, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, resp)
}

func (h *ChannelHandler) CreateChannel(ctx *gin.Context) {
	req := new(apiV1.CreateChannelRequest)
	if err := ctx.ShouldBind(req); err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	resp, err := h.svc.CreateChannel(ctx, req)
	if err != nil {
		apiV1.HandleError(ctx, 0, err, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, strconv.FormatUint(resp, 10))
}

func (h *ChannelHandler) DeleteChannel(ctx *gin.Context) {
	channelId := ctx.Param("channelId")
	if channelId == "" {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, "channelId is required")
		return
	}
	channelIdUint, err := strconv.ParseUint(channelId, 10, 64)
	if err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, "channelId is invalid")
		return
	}
	err = h.svc.DeleteChannel(ctx, channelIdUint)
	if err != nil {
		apiV1.HandleError(ctx, 0, err, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, nil)
}

func (h *ChannelHandler) GetChannel(ctx *gin.Context) {
	channelId := ctx.Param("channelId")
	if channelId == "" {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, "channelId is required")
		return
	}
	channelIdUint, err := strconv.ParseUint(channelId, 10, 64)
	if err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, "channelId is invalid")
		return
	}
	resp, err := h.svc.GetChannel(ctx, channelIdUint)
	if err != nil {
		apiV1.HandleError(ctx, 0, err, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, resp)
}

func (h *ChannelHandler) UpdateChannel(ctx *gin.Context) {
	var req apiV1.UpdateChannelRequest
	if err := ctx.ShouldBind(&req); err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	channelId := ctx.Param("channelId")
	if channelId == "" {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, "channelId is required")
		return
	}
	channelIdUint, err := strconv.ParseUint(ctx.Param("channelId"), 10, 64)
	if err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, "channelId is invalid")
		return
	}
	err = h.svc.UpdateChannel(ctx, channelIdUint, &req)
	if err != nil {
		apiV1.HandleError(ctx, 0, err, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, nil)
}

func (h *ChannelHandler) UpdateChannelStatus(ctx *gin.Context) {
	var req apiV1.UpdateChannelRequest
	if err := ctx.ShouldBind(&req); err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	channelId := ctx.Param("channelId")
	if channelId == "" {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, "channelId is required")
		return
	}
	channelIdUint, err := strconv.ParseUint(ctx.Param("channelId"), 10, 64)
	if err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, "channelId is invalid")
		return
	}
	err = h.svc.UpdateChannelStatus(ctx, channelIdUint, req.Status)
	if err != nil {
		apiV1.HandleError(ctx, 0, err, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, nil)
}

func (h *ChannelHandler) CheckModel(ctx *gin.Context) {
	channelId := ctx.Param("channelId")
	if channelId == "" {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, "channelId is required")
		return
	}
	channelIdUint, err := strconv.ParseUint(channelId, 10, 64)
	if err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, "channelId is invalid")
		return
	}

	var req apiV1.CheckModelRequest
	if err = ctx.ShouldBind(&req); err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}

	modelId := req.ModelName
	resp, err := h.checkSvc.CheckModel2(ctx, channelIdUint, modelId)
	if err != nil {
		apiV1.HandleError(ctx, 0, err, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, resp)
}
