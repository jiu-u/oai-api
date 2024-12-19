package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jiu-u/oai-api/internal/model"
	"github.com/jiu-u/oai-api/internal/repository"
	"github.com/spf13/viper"
	"math/rand"
	"sync"
)

type ProviderConf struct {
	ProviderName string
	ProviderType string
	ProviderId   uint64
	EndPoint     string
	APIKey       string
	ModelUID     uint64
	ModelId      string
}

type LoadBalanceService interface {
	NextProvider(ctx context.Context, modelId string, statusName string) (*ProviderConf, error)
	AddProvider(ctx context.Context, provider *model.Provider) error
	RemoveProvider(ctx context.Context, id uint64) error
	ChangeModelMapping(ctx context.Context, modelMapping map[string][]string)
}

type loadBalanceService struct {
	*Service
	mu           *sync.RWMutex
	providerMap  map[uint64]*model.Provider
	providerRepo repository.ProviderRepo
	modelRepo    repository.ModelRepo
	modelMapping map[string][]string
}

func NewLoadBalanceService(
	s *Service,
	providerRepo repository.ProviderRepo,
	modelRepo repository.ModelRepo,
	conf *viper.Viper,
) LoadBalanceService {
	modelMapping := make(map[string][]string)
	if err := viper.UnmarshalKey("routes", &modelMapping); err != nil {
		panic(err)
	}
	list, err := providerRepo.FindAll(context.Background())
	if err != nil {
		panic(err)
	}
	mp := make(map[uint64]*model.Provider)
	for _, provider := range list {
		mp[provider.Id] = provider
	}
	return &loadBalanceService{
		Service:      s,
		providerMap:  mp,
		providerRepo: providerRepo,
		modelRepo:    modelRepo,
		mu:           &sync.RWMutex{},
		modelMapping: modelMapping,
	}
}

func (l *loadBalanceService) NextProvider(ctx context.Context, modelId, statusName string) (*ProviderConf, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	models, _ := l.modelMapping[modelId]
	models = append(models, modelId)
	result, err := l.modelRepo.FindUsefulModels(ctx, models, statusName)
	if err != nil || len(result) == 0 {
		fmt.Println("err", err)
		return nil, errors.New("no available provider")
	}
	// 随机负载均衡
	// 随机选择一个provider
	totalWeight := 0
	for _, item := range result {
		totalWeight += item.Weight
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
	return &ProviderConf{
		ProviderName: l.providerMap[selected.ProviderId].Name,
		ProviderType: l.providerMap[selected.ProviderId].Type,
		ProviderId:   selected.ProviderId,
		EndPoint:     l.providerMap[selected.ProviderId].EndPoint,
		APIKey:       l.providerMap[selected.ProviderId].APIKey,
		ModelUID:     selected.Id,
		ModelId:      selected.Model,
	}, nil
}

func (l *loadBalanceService) AddProvider(ctx context.Context, provider *model.Provider) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.providerMap[provider.Id] = provider
	return nil
}

func (l *loadBalanceService) RemoveProvider(ctx context.Context, id uint64) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.providerMap, id)
	return nil
}

func (l *loadBalanceService) ChangeModelMapping(ctx context.Context, modelMapping map[string][]string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.modelMapping = modelMapping
}
