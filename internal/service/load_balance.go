package service

import (
	"strconv"
)

//type ChannelModelConf struct {
//	ProviderName   string
//	ProviderType   string
//	ProviderId     uint64
//	EndPoint       string
//	APIKey         string
//	ModelUID       uint64
//	ModelId        string
//	ProviderModels []string
//}

//type LoadBalanceService interface {
//	NextProvider(ctx context.Context, modelId string, statusName string) (*ChannelModelConf, error)
//	AddProvider(ctx context.Context, provider *model.Provider) error
//	RemoveProvider(ctx context.Context, id uint64) error
//	ChangeModelMapping(ctx context.Context, modelMapping map[string][]string)
//	GetModelMappingKeys() []string
//}
//
//type loadBalanceService struct {
//	*Service
//	mu           *sync.RWMutex
//	providerMap  map[uint64]*model.Provider
//	providerRepo repository.ProviderRepo
//	modelRepo    repository.ModelRepo
//	modelMapping map[string][]string
//	once         *sync.Once
//}
//
//func NewLoadBalanceService(
//	s *Service,
//	providerRepo repository.ProviderRepo,
//	modelRepo repository.ModelRepo,
//	cfg *config.Config,
//) LoadBalanceService {
//	modelMapping := cfg.ModelMapping
//	//list, err := providerRepo.FindAll(context.Background())
//	//if err != nil {
//	//	panic(err)
//	//}
//	//mp := make(map[uint64]*model.Provider)
//	//for _, provider := range list {
//	//	mp[provider.Id] = provider
//	//}
//	return &loadBalanceService{
//		Service:      s,
//		providerRepo: providerRepo,
//		modelRepo:    modelRepo,
//		mu:           &sync.RWMutex{},
//		modelMapping: modelMapping,
//		once:         &sync.Once{},
//	}
//}
//
//func (l *loadBalanceService) loadProviderData() {
//	l.mu.Lock()
//	defer l.mu.Unlock()
//	list, err := l.providerRepo.FindAll(context.Background())
//	if err != nil {
//		panic(err)
//	}
//	mp := make(map[uint64]*model.Provider)
//	for _, provider := range list {
//		mp[provider.Id] = provider
//	}
//	l.providerMap = mp
//}
//
//func (l *loadBalanceService) NextProvider(ctx context.Context, modelId, statusName string) (*ChannelModelConf, error) {
//	l.once.Do(l.loadProviderData)
//	l.mu.RLock()
//	defer l.mu.RUnlock()
//	models, _ := l.modelMapping[modelId]
//	models = append(models, modelId)
//	result, err := l.modelRepo.FindUsefulModels(ctx, models, statusName)
//	if err != nil || len(result) == 0 {
//		l.Logger.WithContext(ctx).Warn("no available provider", zap.Error(err))
//		return nil, errors.New("no available provider")
//	}
//	// 随机负载均衡
//	// 随机选择一个provider
//	totalWeight := 0
//	for _, item := range result {
//		totalWeight += item.Weight
//	}
//	if totalWeight == 0 {
//		return nil, errors.New("no available provider")
//	}
//	idx := 0
//	randomWeight := rand.Intn(totalWeight)
//	for i, item := range result {
//		randomWeight -= item.Weight
//		if randomWeight < 0 {
//			idx = i
//			break
//		}
//	}
//	selected := result[idx]
//
//	str := fmt.Sprintf("请求||请求模型:%s\t\t 实际模型:%s\t\t 服务商名称：%s\t\t ID:%v\t\t Content:%s\n", modelId, selected.Model, l.providerMap[selected.ProviderId].Name, selected.ProviderId, GetKeyId(l.providerMap[selected.ProviderId].APIKey, selected.ProviderId))
//	l.Logger.WithContext(ctx).Info(str)
//	return &ChannelModelConf{
//		ProviderName: l.providerMap[selected.ProviderId].Name,
//		ProviderType: l.providerMap[selected.ProviderId].Type,
//		ProviderId:   selected.ProviderId,
//		EndPoint:     l.providerMap[selected.ProviderId].EndPoint,
//		APIKey:       l.providerMap[selected.ProviderId].APIKey,
//		ModelUID:     selected.Id,
//		ModelId:      selected.Model,
//	}, nil
//}
//
//func (l *loadBalanceService) AddProvider(ctx context.Context, provider *model.Provider) error {
//	l.mu.Lock()
//	defer l.mu.Unlock()
//	l.providerMap[provider.Id] = provider
//	return nil
//}
//
//func (l *loadBalanceService) RemoveProvider(ctx context.Context, id uint64) error {
//	l.mu.Lock()
//	defer l.mu.Unlock()
//	delete(l.providerMap, id)
//	return nil
//}
//
//func (l *loadBalanceService) ChangeModelMapping(ctx context.Context, modelMapping map[string][]string) {
//	l.mu.Lock()
//	defer l.mu.Unlock()
//	l.modelMapping = modelMapping
//}
//
//func (l *loadBalanceService) GetModelMappingKeys() []string {
//	l.mu.RLock()
//	defer l.mu.RUnlock()
//	keys := make([]string, 0, len(l.modelMapping))
//	for k := range l.modelMapping {
//		keys = append(keys, k)
//	}
//	return keys
//}

func GetKeyId(key string, id uint64) string {
	keyId := "不足8位，id->" + strconv.FormatUint(id, 10)
	if len(key) >= 8 {
		keyId = key[len(key)-8:]
	}
	return keyId
}
