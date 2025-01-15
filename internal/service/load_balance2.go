package service

import (
	"context"
	"errors"
	"github.com/jiu-u/oai-api/internal/dto"
	"github.com/jiu-u/oai-api/internal/model"
	"github.com/jiu-u/oai-api/internal/repository"
	"go.uber.org/zap"
	"math/rand"
	"sync"
	"time"
)

type LoadBalanceServiceBeta interface {
	AddChannel(ctx context.Context, channel *model.Channel) error
	RemoveChannel(ctx context.Context, id uint64) error
	NextChannel(ctx context.Context, modelId string) (*dto.ChannelModelConf, error)
	SuccessCb(ctx context.Context, modelRecordId uint64) error
	FailCb(ctx context.Context, modelRecordId uint64) error
	ChangeModelMapping(ctx context.Context, modelMapping map[string][]string)
	RecoverChannelModels(ctx context.Context) error
	GetModelMappingKeys() []string
}

func NewLoadBalanceServiceBeta(service *Service, channelRepo repository.ChannelRepository, channelModelRepo repository.ChannelModelRepository) LoadBalanceServiceBeta {
	return &loadBalanceServiceBeta{
		Service:          service,
		channelRepo:      channelRepo,
		channelModelRepo: channelModelRepo,
		ChannelMap:       make(map[uint64]*model.Channel),
		ModelMapping:     make(map[string][]string),
		RecoverInterval:  5 * time.Minute,
		once:             &sync.Once{},
		mu:               &sync.RWMutex{},
	}
}

type loadBalanceServiceBeta struct {
	*Service
	mu               *sync.RWMutex
	channelRepo      repository.ChannelRepository
	channelModelRepo repository.ChannelModelRepository
	ChannelMap       map[uint64]*model.Channel
	ModelMapping     map[string][]string
	RecoverInterval  time.Duration
	once             *sync.Once
}

func (s *loadBalanceServiceBeta) AddChannel(ctx context.Context, channel *model.Channel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ChannelMap[channel.Id] = channel
	return nil
}

func (s *loadBalanceServiceBeta) RemoveChannel(ctx context.Context, id uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.ChannelMap, id)
	return nil
}

func (s *loadBalanceServiceBeta) loadProviderData() {
	s.mu.Lock()
	defer s.mu.Unlock()
	list, err := s.channelRepo.FindAllChannels(context.Background())
	if err != nil {
		panic(err)
	}
	mp := make(map[uint64]*model.Channel)
	for _, channel := range list {
		mp[channel.Id] = channel
	}
	s.ChannelMap = mp
}

func (s *loadBalanceServiceBeta) NextChannel(ctx context.Context, modelId string) (*dto.ChannelModelConf, error) {
	s.once.Do(s.loadProviderData)
	s.mu.RLock()
	defer s.mu.RUnlock()
	models, _ := s.ModelMapping[modelId]
	models = append(models, modelId)
	result, err := s.channelModelRepo.FindUsefulChannelModels(ctx, models)
	if err != nil || len(result) == 0 {
		s.Logger.WithContext(ctx).Warn("no available provider", zap.Error(err))
		return nil, errors.New("no available provider")
	}
	// 随机负载均衡
	// 随机选择一个channel
	totalWeight := 0
	for _, item := range result {
		totalWeight += item.Weight
	}
	if totalWeight == 0 {
		return nil, errors.New("no available provider")
	}
	idx := 0
	randomWeight := rand.Intn(totalWeight)
	for i, item := range result {
		randomWeight -= item.Weight
		if randomWeight < 0 {
			idx = i
			break
		}
	}
	selected := result[idx]
	channel := s.ChannelMap[selected.ChannelId]
	if channel == nil {
		return nil, errors.New("channel is nil")
	}
	return &dto.ChannelModelConf{
		ChannelId:       selected.ChannelId,
		ChannelName:     channel.Name,
		ChannelType:     channel.Type,
		ChannelKey:      channel.APIKey,
		ChannelEndPoint: channel.EndPoint,
		ModelRecordId:   selected.Id,
		ModelKey:        selected.ModelKey,
		ModelId:         selected.ModelKey,
		Weight:          selected.Weight,
	}, nil

}

func (s *loadBalanceServiceBeta) SuccessCb(ctx context.Context, modelRecordId uint64) error {
	return s.channelModelRepo.InCrChannelModelWeight(ctx, modelRecordId)
}

func (s *loadBalanceServiceBeta) FailCb(ctx context.Context, modelRecordId uint64) error {
	return s.channelModelRepo.DecrChannelModelWeight(ctx, modelRecordId)
}

func (s *loadBalanceServiceBeta) ChangeModelMapping(ctx context.Context, modelMapping map[string][]string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ModelMapping = modelMapping
}

func (s *loadBalanceServiceBeta) RecoverChannelModels(ctx context.Context) error {
	return s.channelModelRepo.RestoreChannelModel(ctx)
}

func (s *loadBalanceServiceBeta) GetModelMappingKeys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.ModelMapping))
	for k := range s.ModelMapping {
		keys = append(keys, k)
	}
	return keys
}
