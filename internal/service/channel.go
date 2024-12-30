package service

import (
	"context"
	"fmt"
	v1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/model"
	"github.com/jiu-u/oai-api/internal/repository"
	"github.com/jiu-u/oai-api/pkg/encrypte"
	"go.uber.org/zap"
	"time"
)

type ChannelService interface {
	CreateChannel(ctx context.Context, provider *v1.CreateProviderRequest) (uint64, error)
}

func NewChannelService(
	srv *Service,
	repo repository.ChannelRepository,
	channelModelRepo repository.ChannelModelRepository,
) ChannelService {
	return &channelService{
		Service:          srv,
		repo:             repo,
		channelModelRepo: channelModelRepo,
	}
}

type channelService struct {
	*Service
	repo             repository.ChannelRepository
	channelModelRepo repository.ChannelModelRepository
}

func (s *channelService) CreateChannel(ctx context.Context, req *v1.CreateProviderRequest) (uint64, error) {
	id := s.Sid.GenUint64()
	format := "%s@%s@%s"
	hashId := encrypte.Sha256Encode(fmt.Sprintf(format, req.Type, req.EndPoint, req.APIKey))
	channel := &model.Channel{
		Name:     req.Name,
		Type:     req.Type,
		EndPoint: req.EndPoint,
		APIKey:   req.APIKey,
		HashId:   hashId,
	}
	channel.Id = id
	fmt.Println("id__>", id)
	err := s.Tm.Transaction(ctx, func(ctx context.Context) error {
		err := s.repo.CreateChannelIfNotExists(ctx, channel)
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
	return id, nil
}
