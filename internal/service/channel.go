package service

import (
	"context"
	"fmt"
	v1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/dto/query"
	"github.com/jiu-u/oai-api/internal/model"
	"github.com/jiu-u/oai-api/internal/repository"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type ChannelService interface {
	CreateChannel(ctx context.Context, provider *v1.CreateChannelRequest) (uint64, error)
	GetChannels(ctx context.Context, req *query.ChannelQueryRequest) (*v1.ChannelListResponse, error)
	DeleteChannel(ctx context.Context, channelId uint64) error
	GetChannel(ctx context.Context, channelId uint64) (*v1.ChannelResponse, error)
	UpdateChannel(ctx context.Context, channelId uint64, req *v1.UpdateChannelRequest) error
	UpdateChannelStatus(ctx context.Context, channelId uint64, status int8) error
}

func NewChannelService(
	srv *Service,
	repo repository.ChannelRepository,
	channelModelRepo repository.ChannelModelRepository,
	loadSvc LoadBalanceServiceBeta,
) ChannelService {
	return &channelService{
		Service:          srv,
		repo:             repo,
		channelModelRepo: channelModelRepo,
		loadSvc:          loadSvc,
	}
}

type channelService struct {
	*Service
	repo             repository.ChannelRepository
	channelModelRepo repository.ChannelModelRepository
	loadSvc          LoadBalanceServiceBeta
}

func (s *channelService) CreateChannel(ctx context.Context, req *v1.CreateChannelRequest) (uint64, error) {
	id := s.Sid.GenUint64()
	channel := &model.Channel{
		Name:     req.Name,
		Type:     req.Type,
		EndPoint: req.EndPoint,
		APIKey:   req.APIKey,
	}
	channel.GenerateHashId()
	channel.Id = id
	err := s.Tm.Transaction(ctx, func(ctx context.Context) error {
		err := s.repo.CreateChannel(ctx, channel)
		if err != nil {
			return fmt.Errorf("create channel failed: %s", err)
		}
		for _, modelId := range req.Models {
			newId := s.Sid.GenUint64()
			newModel := &model.ChannelModel{
				ChannelId:     channel.Id,
				ModelKey:      modelId,
				Weight:        req.Weight,
				LastCheckTime: time.Now(),
				ErrorCount:    0,
				TotalCount:    0,
			}
			newModel.Id = newId
			err = s.channelModelRepo.CreateChannelModelIfNotExists(ctx, newModel)
			if err != nil {
				s.Logger.WithContext(ctx).Warn("create channel model failed", zap.Error(err))
				continue
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	_ = s.loadSvc.AddChannel(ctx, channel)
	return id, nil
}

func (s *channelService) GetChannels(ctx context.Context, req *query.ChannelQueryRequest) (*v1.ChannelListResponse, error) {
	resp := new(v1.ChannelListResponse)
	channels, total, err := s.repo.FindAllChannelsByCondition(ctx, req)
	if err != nil {
		return nil, err
	}
	resp.Total = total
	resp.Page = int64(req.Page)
	resp.PageSize = int64(req.PageSize)
	resp.List = make([]v1.ChannelResponse, len(channels))
	for idx, channel := range channels {
		resp.List[idx] = v1.ChannelResponse{
			Id:       strconv.FormatUint(channel.Id, 10),
			Name:     channel.Name,
			Type:     channel.Type,
			Balance:  channel.Balance,
			EndPoint: channel.EndPoint,
			APIKey:   channel.APIKey,
			Models:   make([]string, len(channel.Models)),
			Status:   channel.Status,
		}
		for jdx, modelX := range channel.Models {
			resp.List[idx].Models[jdx] = modelX.ModelKey
		}
	}
	return resp, nil
}

func (s *channelService) DeleteChannel(ctx context.Context, channelId uint64) error {
	var err error
	err = s.Tm.Transaction(ctx, func(ctx context.Context) error {
		err = s.repo.DeleteChannelByID(ctx, channelId)
		if err != nil {
			return err
		}
		err = s.channelModelRepo.DeleteChannelModelByChannelId(ctx, channelId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	_ = s.loadSvc.RemoveChannel(ctx, channelId)
	return nil
}

func (s *channelService) GetChannel(ctx context.Context, channelId uint64) (*v1.ChannelResponse, error) {
	var resp v1.ChannelResponse
	var err error
	err = s.Tm.Transaction(ctx, func(ctx context.Context) error {
		channel, err := s.repo.FindChannelById(ctx, channelId)
		if err != nil {
			return err
		}
		resp.Models = make([]string, len(channel.Models))
		for idx, modelX := range channel.Models {
			resp.Models[idx] = modelX.ModelKey
		}
		resp.Id = strconv.FormatUint(channel.Id, 10)
		resp.Status = channel.Status
		resp.Name = channel.Name
		resp.Type = channel.Type
		resp.Balance = channel.Balance
		resp.EndPoint = channel.EndPoint
		resp.APIKey = channel.APIKey
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *channelService) UpdateChannel(ctx context.Context, channelId uint64, req *v1.UpdateChannelRequest) error {
	var err error
	var channelX *model.Channel
	err = s.Tm.Transaction(ctx, func(ctx context.Context) error {
		channelX = &model.Channel{
			Name:     req.Name,
			Type:     req.Type,
			EndPoint: req.EndPoint,
			APIKey:   req.APIKey,
			Status:   req.Status,
			Models:   nil,
		}
		channelX.Id = channelId
		err = s.repo.UpdateChannel(ctx, channelX)
		if err != nil {
			return err
		}
		channelModels := make([]*model.ChannelModel, len(req.Models))
		for idx, modelKey := range req.Models {
			id := s.Sid.GenUint64()
			channelModels[idx] = &model.ChannelModel{
				ChannelId:     channelId,
				ModelKey:      modelKey,
				Weight:        10,
				LastCheckTime: time.Now(),
				ErrorCount:    0,
				TotalCount:    0,
			}
			channelModels[idx].Id = id
		}
		if len(channelModels) > 0 {
			err = s.channelModelRepo.ResetChannelModels(ctx, channelId, channelModels)
			if err != nil {
				return err
			}
		}
		if req.Status > 0 && req.Status < 3 {
			err = s.channelModelRepo.UpdateChannelModelsHardStatus(ctx, channelId, req.Status)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if channelX.Status == 1 {
		_ = s.loadSvc.AddChannel(ctx, channelX)
	} else if channelX.Status == 2 {
		_ = s.loadSvc.RemoveChannel(ctx, channelX.Id)
	}
	return err
}

func (s *channelService) UpdateChannelStatus(ctx context.Context, channelId uint64, status int8) error {
	return s.UpdateChannel(ctx, channelId, &v1.UpdateChannelRequest{
		Status: status,
	})
}
