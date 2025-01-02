package service

import (
	"context"
	"errors"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/repository"
	adapterV1 "github.com/jiu-u/oai-api/pkg/adapter/api/v1"
	"github.com/jiu-u/oai-api/pkg/adapter/provider"
	"github.com/jiu-u/oai-api/pkg/array"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

type RelayType int

const (
	RelayChat RelayType = iota
	RelayCompletion
	RelayEmbedding
	RelaySpeech
	RelayTranscriptions
	RelayTranslations
	RelayImage
	RelayImageEdit
	RelayImageVariations
	RelayChatByBytes
	RelayCompletionByBytes
	RelayEmbeddingByBytes
	RelaySpeechByBytes
	RelayImageByBytes
)

type OaiService interface {
	RelayRequest(ctx context.Context, req any, modelId string, relayType RelayType) (io.ReadCloser, http.Header, error)
	ChatCompletions(ctx context.Context, req *adapterV1.ChatCompletionRequest) (io.ReadCloser, http.Header, error)
	ChatCompletionsByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error)
	Models(ctx context.Context) (*apiV1.ModelResponse, error)
	Completions(ctx context.Context, req *adapterV1.CompletionsRequest) (io.ReadCloser, http.Header, error)
	CompletionsByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error)
	Embeddings(ctx context.Context, req *adapterV1.EmbeddingRequest) (io.ReadCloser, http.Header, error)
	EmbeddingsByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error)
	CreateSpeech(ctx context.Context, req *adapterV1.SpeechRequest) (io.ReadCloser, http.Header, error)
	CreateSpeechByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error)
	Transcriptions(ctx context.Context, req *adapterV1.TranscriptionRequest) (io.ReadCloser, http.Header, error)
	Translations(ctx context.Context, req *adapterV1.TranslationRequest) (io.ReadCloser, http.Header, error)
	CreateImage(ctx context.Context, req *adapterV1.CreateImageRequest) (io.ReadCloser, http.Header, error)
	CreateImageByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error)
	CreateImageEdit(ctx context.Context, req *adapterV1.EditImageRequest) (io.ReadCloser, http.Header, error)
	ImageVariations(ctx context.Context, req *adapterV1.CreateImageVariationRequest) (io.ReadCloser, http.Header, error)
}

func NewOAIProvider(conf *ChannelModelConf) (provider.Provider, error) {
	aConf := provider.Config{
		Type:     conf.ChannelType,
		EndPoint: conf.ChannelEndPoint,
		APIKey:   conf.ChannelKey,
	}
	switch aConf.Type {
	case "openai":
		return provider.NewOpenAIProvider(aConf), nil
	case "oaiNoFetchModel":
		return provider.NewOaiNoFetchModelProvider(aConf, nil), nil
	case "siliconflow":
		return provider.NewSiliconFlowProvider(aConf), nil
	default:
		return nil, errors.New("invalid provider type")
	}
}

func NewOaiService(
	svc *Service,
	load LoadBalanceServiceBeta,
	reqLogSvc RequestLogService,
	channelModelRepo repository.ChannelModelRepository,
) OaiService {
	return &oaiService{
		Service:          svc,
		load:             load,
		N:                3,
		reqLogSvc:        reqLogSvc,
		channelModelRepo: channelModelRepo,
	}
}

type oaiService struct {
	*Service
	load             LoadBalanceServiceBeta
	N                int
	channelModelRepo repository.ChannelModelRepository
	reqLogSvc        RequestLogService
}

func (s *oaiService) GoLogReq(ctx context.Context, modelId string, status int8) {
	go s.LogReq(ctx, modelId, status)
}

func (s *oaiService) LogReq(ctx context.Context, modelId string, status int8) {
	_, ok := ctx.(*gin.Context)
	if !ok {
		return
	}
	apiKey, err := GetApiKey(ctx.(*gin.Context))
	if err != nil {
		return
	}
	req := RequestLogReq{
		Key:    apiKey,
		Status: status,
		Model:  modelId,
		Ip:     ctx.(*gin.Context).ClientIP(),
	}
	err = s.reqLogSvc.CreateRequestLog(ctx, &req)
	if err != nil {
		s.Logger.Warn("创建请求日志失败", zap.Error(err))
	}
}

func GetApiKey(ctx *gin.Context) (string, error) {
	apiKey, exist := ctx.Get("apiKey")
	if apiKey == "" || !exist {
		return "", errors.New("api key is empty")
	}
	return apiKey.(string), nil
}

func (s *oaiService) SuccessCb(ctx context.Context, modelRecordId uint64) {
	go func() {
		err := s.load.SuccessCb(ctx, modelRecordId)
		if err != nil {
			s.Logger.Warn("successCb失败", zap.Error(err))
		}
	}()
}

func (s *oaiService) FailCb(ctx context.Context, modelRecordId uint64) {
	err := s.load.FailCb(ctx, modelRecordId)
	if err != nil {
		s.Logger.Warn("failCb失败", zap.Error(err))
	}
}

func (s *oaiService) RelayRequest(ctx context.Context, req any, modelId string, relayType RelayType) (io.ReadCloser, http.Header, error) {
	reqModelId := modelId
	if reqModelId == "" {
		return nil, nil, errors.New("modelId is empty")
	}
	logger := s.Logger.WithContext(ctx)
	for i := range s.N {
		zapLogger := logger.With(
			zap.String("reqModelId", reqModelId),
			zap.String("relayType", strconv.Itoa(int(relayType))),
			zap.Int("loop_times", i),
		)
		conf, err := s.load.NextChannel(ctx, reqModelId)
		if err != nil {
			zapLogger.Warn("获取provider失败", zap.Error(err))
			continue
		}
		newProvider, err := NewOAIProvider(conf)
		if err != nil {
			zapLogger.Warn("创建provider失败", zap.Error(err))
			continue
		}
		resp, respHeader, err := s.DoRelayRequest(ctx, req, conf.ModelKey, relayType, newProvider)
		if err == nil {
			zapLogger.Info("获取response成功", zap.Error(err))
			s.SuccessCb(ctx, conf.ModelRecordId)
			s.GoLogReq(ctx, conf.ModelKey, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	return nil, nil, errors.New("all provider failed.please try again later")
}

func (s *oaiService) DoRelayRequest(ctx context.Context, reqBody any, modelId string, relayType RelayType, p provider.Provider) (io.ReadCloser, http.Header, error) {
	switch relayType {
	case RelayChat:
		req, ok := reqBody.(*adapterV1.ChatCompletionRequest)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req.Model = modelId
		return p.ChatCompletions(ctx, req)
	case RelayChatByBytes:
		req, ok := reqBody.([]byte)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req, err := changeBytesModelId(req, modelId)
		if err != nil {
			return nil, nil, err
		}
		return p.ChatCompletionsByBytes(ctx, req)
	case RelayCompletion:
		req, ok := reqBody.(*adapterV1.CompletionsRequest)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req.Model = modelId
		return p.Completions(ctx, req)
	case RelayCompletionByBytes:
		req, ok := reqBody.([]byte)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req, err := changeBytesModelId(req, modelId)
		if err != nil {
			return nil, nil, err
		}
		return p.CompletionsByBytes(ctx, req)
	case RelayEmbedding:
		req, ok := reqBody.(*adapterV1.EmbeddingRequest)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req.Model = modelId
		return p.Embeddings(ctx, req)
	case RelayEmbeddingByBytes:
		req, ok := reqBody.([]byte)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req, err := changeBytesModelId(req, modelId)
		if err != nil {
			return nil, nil, err
		}
		return p.EmbeddingsByBytes(ctx, req)
	case RelaySpeech:
		req, ok := reqBody.(*adapterV1.SpeechRequest)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req.Model = modelId
		return p.CreateSpeech(ctx, req)
	case RelaySpeechByBytes:
		req, ok := reqBody.([]byte)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req, err := changeBytesModelId(req, modelId)
		if err != nil {
			return nil, nil, err
		}
		return p.CreateSpeechByBytes(ctx, req)
	case RelayTranscriptions:
		req, ok := reqBody.(*adapterV1.TranscriptionRequest)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req.Model = modelId
		return p.Transcriptions(ctx, req)
	case RelayTranslations:
		req, ok := reqBody.(*adapterV1.TranslationRequest)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req.Model = modelId
		return p.Translations(ctx, req)
	case RelayImage:
		req, ok := reqBody.(*adapterV1.CreateImageRequest)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req.Model = modelId
		return p.CreateImage(ctx, req)
	case RelayImageByBytes:
		req, ok := reqBody.([]byte)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req, err := changeBytesModelId(req, modelId)
		if err != nil {
			return nil, nil, err
		}
		return p.CreateImageByBytes(ctx, req)
	case RelayImageEdit:
		req, ok := reqBody.(*adapterV1.EditImageRequest)
		if !ok {
			return nil, nil, errors.New("invalid request body")

		}
		req.Model = modelId
		return p.CreateImageEdit(ctx, req)
	case RelayImageVariations:
		req, ok := reqBody.(*adapterV1.CreateImageVariationRequest)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req.Model = modelId
		return p.ImageVariations(ctx, req)
	}
	return nil, nil, errors.New("invalid relay type")
}

func (s *oaiService) ChatCompletions2Archive(ctx context.Context, req *adapterV1.ChatCompletionRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, reqModelId)
		if err != nil {
			s.Logger.Warn("获取provider失败", zap.String("modelId", reqModelId), zap.Error(err))
			continue
		}
		newProvider, err := NewOAIProvider(conf)
		if err != nil {
			s.Logger.Warn("创建provider失败", zap.String("modelId", reqModelId), zap.Error(err))
			continue
		}
		req.Model = conf.ModelKey
		resp, respHeader, err := newProvider.ChatCompletions(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			s.GoLogReq(ctx, conf.ModelKey, 1)
			return resp, respHeader, nil
		}
		detail, _ := io.ReadAll(resp)
		s.Logger.Warn("获取response失败",
			zap.String("modelId", conf.ModelKey),
			zap.String("channel", strconv.FormatUint(conf.ChannelId, 10)),
			zap.String("providerType", conf.ChannelType),
			zap.String("providerName", conf.ChannelName),
			zap.String("detail", string(detail)),
			zap.Error(err))
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
		s.Logger.Warn("更新状态失败",
			zap.String("modelId", conf.ModelKey),
			zap.String("provider", strconv.FormatUint(conf.ChannelId, 10)),
			zap.Error(err))
	}
	s.GoLogReq(ctx, reqModelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) ChatCompletionsByBytes2Archive(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, modelId)
		if err != nil {
			continue
		}
		newProvider, err := NewOAIProvider(conf)
		if err != nil {
			continue
		}
		req, err = changeBytesModelId(req, conf.ModelId)
		resp, respHeader, err := newProvider.ChatCompletionsByBytes(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	s.GoLogReq(ctx, modelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) Models(ctx context.Context) (*apiV1.ModelResponse, error) {
	modelIds, err := s.channelModelRepo.FindAllChannelModelIds(ctx)
	if err != nil {
		return nil, err
	}
	resp := new(apiV1.ModelResponse)
	modelSet := make(map[string]struct{}, len(modelIds))
	resp.Object = "list"
	resp.Data = array.Map(modelIds, func(modelId string) adapterV1.Model {
		modelSet[modelId] = struct{}{}
		return adapterV1.Model{
			ID:      modelId,
			Object:  "model",
			Created: 0,
		}
	})
	list := s.load.GetModelMappingKeys()
	for _, modelId := range list {
		if _, ok := modelSet[modelId]; !ok {
			resp.Data = append(resp.Data, adapterV1.Model{
				ID:      modelId,
				Object:  "model",
				Created: 0,
			})
		}
	}
	return resp, nil
}

func (s *oaiService) Completions2Archive(ctx context.Context, req *adapterV1.CompletionsRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, reqModelId)
		if err != nil {
			continue
		}
		newProvider, err := NewOAIProvider(conf)
		if err != nil {
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := newProvider.Completions(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	s.GoLogReq(ctx, reqModelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) CompletionsByBytes2Archive(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, modelId)
		if err != nil {
			continue
		}
		newProvider, err := NewOAIProvider(conf)
		if err != nil {
			continue
		}
		req, err = changeBytesModelId(req, conf.ModelId)
		resp, respHeader, err := newProvider.CompletionsByBytes(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	s.GoLogReq(ctx, modelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) Embeddings2Archive(ctx context.Context, req *adapterV1.EmbeddingRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, reqModelId)
		if err != nil {
			continue
		}
		newProvider, err := NewOAIProvider(conf)
		if err != nil {
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := newProvider.Embeddings(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	s.GoLogReq(ctx, reqModelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) EmbeddingsByBytes2Archive(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, modelId)
		if err != nil {
			continue
		}
		newProvider, err := NewOAIProvider(conf)
		if err != nil {
			continue
		}
		req, err = changeBytesModelId(req, conf.ModelId)
		resp, respHeader, err := newProvider.EmbeddingsByBytes(ctx, req)
		if err == nil {
			s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	s.GoLogReq(ctx, modelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) CreateSpeech2Archive(ctx context.Context, req *adapterV1.SpeechRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, reqModelId)
		if err != nil {
			continue
		}
		newProvider, err := NewOAIProvider(conf)
		if err != nil {
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := newProvider.CreateSpeech(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	s.GoLogReq(ctx, reqModelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) CreateSpeechByBytes2Archive(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, modelId)
		if err != nil {
			continue
		}
		newProvider, err := NewOAIProvider(conf)
		if err != nil {
			continue
		}
		req, err = changeBytesModelId(req, conf.ModelId)
		resp, respHeader, err := newProvider.CreateSpeechByBytes(ctx, req)
		if err == nil {
			s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	s.GoLogReq(ctx, modelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) Transcriptions2Archive(ctx context.Context, req *adapterV1.TranscriptionRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, reqModelId)
		if err != nil {
			continue
		}
		newProvider, err := NewOAIProvider(conf)
		if err != nil {
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := newProvider.Transcriptions(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	s.GoLogReq(ctx, reqModelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) Translations2Archive(ctx context.Context, req *adapterV1.TranslationRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, reqModelId)
		if err != nil {
			continue
		}
		newProvider, err := NewOAIProvider(conf)
		if err != nil {
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := newProvider.Translations(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	s.GoLogReq(ctx, reqModelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) CreateImage2Archive(ctx context.Context, req *adapterV1.CreateImageRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, reqModelId)
		if err != nil {
			continue
		}
		newProvider, err := NewOAIProvider(conf)
		if err != nil {
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := newProvider.CreateImage(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	s.GoLogReq(ctx, reqModelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) CreateImageByBytes2Archive(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, modelId)
		if err != nil {
			continue
		}
		newProvider, err := NewOAIProvider(conf)
		if err != nil {
			continue
		}
		req, err = changeBytesModelId(req, conf.ModelId)
		resp, respHeader, err := newProvider.CreateImageByBytes(ctx, req)
		if err == nil {
			s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	s.GoLogReq(ctx, modelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) CreateImageEdit2Archive(ctx context.Context, req *adapterV1.EditImageRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, reqModelId)
		if err != nil {
			continue
		}
		newProvider, err := NewOAIProvider(conf)
		if err != nil {
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := newProvider.CreateImageEdit(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		s.FailCb(ctx, conf.ModelRecordId)
	}
	s.GoLogReq(ctx, reqModelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) ImageVariations2Archive(ctx context.Context, req *adapterV1.CreateImageVariationRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, reqModelId)
		if err != nil {
			continue
		}
		newProvider, err := NewOAIProvider(conf)
		if err != nil {
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := newProvider.ImageVariations(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		s.FailCb(ctx, conf.ModelRecordId)
	}
	s.GoLogReq(ctx, reqModelId, 2)
	return nil, nil, errors.New("no service available")
}

func changeModelId(req *adapterV1.ChatCompletionRequest, newModelId string) {
	req.Model = newModelId
}

func changeBytesModelId(bodyBytes []byte, newModelId string) ([]byte, error) {
	var result map[string]any
	err := sonic.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}
	result["model"] = newModelId
	bytes, err := sonic.Marshal(result)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
