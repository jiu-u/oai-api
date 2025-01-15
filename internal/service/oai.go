package service

import (
	"context"
	"errors"
	"github.com/bytedance/sonic"
	adapter "github.com/jiu-u/oai-adapter"
	adapterApi "github.com/jiu-u/oai-adapter/api"
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/dto"
	"github.com/jiu-u/oai-api/internal/repository"
	adapterV1 "github.com/jiu-u/oai-api/pkg/adapter/api/v1"
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
	ChatCompletions(ctx context.Context, req *adapterApi.ChatRequest) (io.ReadCloser, http.Header, error)
	ChatCompletionsByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error)
	Models(ctx context.Context) (*apiV1.ModelResponse, error)
	Completions(ctx context.Context, req *adapterApi.CompletionsRequest) (io.ReadCloser, http.Header, error)
	CompletionsByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error)
	Embeddings(ctx context.Context, req *adapterApi.EmbeddingRequest) (io.ReadCloser, http.Header, error)
	EmbeddingsByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error)
	CreateSpeech(ctx context.Context, req *adapterApi.SpeechRequest) (io.ReadCloser, http.Header, error)
	CreateSpeechByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error)
	Transcriptions(ctx context.Context, req *adapterApi.TranscriptionRequest) (io.ReadCloser, http.Header, error)
	Translations(ctx context.Context, req *adapterApi.TranslationRequest) (io.ReadCloser, http.Header, error)
	CreateImage(ctx context.Context, req *adapterApi.CreateImageRequest) (io.ReadCloser, http.Header, error)
	CreateImageByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error)
	CreateImageEdit(ctx context.Context, req *adapterApi.EditImageRequest) (io.ReadCloser, http.Header, error)
	ImageVariations(ctx context.Context, req *adapterApi.CreateImageVariationRequest) (io.ReadCloser, http.Header, error)
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

var typeMp = map[string]adapter.AdapterType{
	"openai":          adapter.OpenAI,
	"oaiNoModels":     adapter.OAINoModels,
	"gemini":          adapter.Gemini,
	"siliconflow":     adapter.SiliconFlow,
	"siliconflowFree": adapter.SiliconFlowFree,
}

func NewOAIAdapter(conf *dto.ChannelModelConf) (adapter.Adapter, error) {
	if _, exist := typeMp[conf.ChannelType]; !exist {
		return nil, errors.New("invalid provider type")
	}
	cfg := &adapter.AdapterConfig{
		AdapterType:  typeMp[conf.ChannelType],
		ApiKey:       conf.ChannelKey,
		EndPoint:     conf.ChannelEndPoint,
		ManualModels: nil,
		ProxyURL:     nil,
	}
	adapter2 := adapter.NewAdapter(cfg)
	return adapter2, nil
}

func (s *oaiService) GoLogReq(ctx context.Context, trace *RequestLogReq) {
	go s.LogReq(ctx, trace)
}

func (s *oaiService) LogReq(ctx context.Context, trace *RequestLogReq) {
	apiKey, err := GetApiKey(ctx)
	if err != nil {
		return
	}
	req := trace
	req.Key = apiKey
	req.Ip = GetClientIp(ctx)
	err = s.reqLogSvc.CreateRequestLog(ctx, req)
	if err != nil {
		s.Logger.Warn("创建请求日志失败", zap.Error(err))
	}
}

func GetApiKey(ctx context.Context) (string, error) {
	apiKey := ctx.Value("apiKey")
	if apiKey == "" {
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
	trace := new(RequestLogReq)
	for i := range s.N {
		trace.RetryTimes = i
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
		trace.Model = conf.ModelKey
		trace.ChannelNames += conf.ChannelName + ","
		trace.ChannelIds += strconv.FormatUint(conf.ModelRecordId, 10) + ","
		adapterX, err := NewOAIAdapter(conf)
		if err != nil {
			zapLogger.Warn("获取provider失败", zap.Error(err))
			continue
		}
		resp, respHeader, err := s.DoRelayRequest(ctx, req, conf.ModelKey, relayType, adapterX)
		if err == nil {
			zapLogger.Info("获取response成功", zap.Error(err))
			s.SuccessCb(ctx, conf.ModelRecordId)
			trace.Status = 1
			s.GoLogReq(ctx, trace)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	trace.Status = 2
	s.GoLogReq(ctx, trace)
	return nil, nil, errors.New("all provider failed.please try again later")
}

func (s *oaiService) DoRelayRequest(ctx context.Context, reqBody any, modelId string, relayType RelayType, ad adapter.Adapter) (io.ReadCloser, http.Header, error) {
	switch relayType {
	case RelayChat:
		req, ok := reqBody.(*adapterApi.ChatRequest)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req.Model = modelId
		return ad.ChatCompletions(ctx, req)
	case RelayChatByBytes:
		req, ok := reqBody.([]byte)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req, err := changeBytesModelId(req, modelId)
		if err != nil {
			return nil, nil, err
		}
		return ad.ChatCompletionsByBytes(ctx, req)
	case RelayCompletion:
		req, ok := reqBody.(*adapterApi.CompletionsRequest)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req.Model = modelId
		return ad.Completions(ctx, req)
	case RelayCompletionByBytes:
		req, ok := reqBody.([]byte)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req, err := changeBytesModelId(req, modelId)
		if err != nil {
			return nil, nil, err
		}
		return ad.CompletionsByBytes(ctx, req)
	case RelayEmbedding:
		req, ok := reqBody.(*adapterApi.EmbeddingRequest)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req.Model = modelId
		return ad.Embeddings(ctx, req)
	case RelayEmbeddingByBytes:
		req, ok := reqBody.([]byte)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req, err := changeBytesModelId(req, modelId)
		if err != nil {
			return nil, nil, err
		}
		return ad.EmbeddingsByBytes(ctx, req)
	case RelaySpeech:
		req, ok := reqBody.(*adapterApi.SpeechRequest)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req.Model = modelId
		return ad.CreateSpeech(ctx, req)
	case RelaySpeechByBytes:
		req, ok := reqBody.([]byte)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req, err := changeBytesModelId(req, modelId)
		if err != nil {
			return nil, nil, err
		}
		return ad.CreateSpeechByBytes(ctx, req)
	case RelayTranscriptions:
		req, ok := reqBody.(*adapterApi.TranscriptionRequest)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req.Model = modelId
		return ad.Transcriptions(ctx, req)
	case RelayTranslations:
		req, ok := reqBody.(*adapterApi.TranslationRequest)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req.Model = modelId
		return ad.Translations(ctx, req)
	case RelayImage:
		req, ok := reqBody.(*adapterApi.CreateImageRequest)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req.Model = modelId
		return ad.CreateImage(ctx, req)
	case RelayImageByBytes:
		req, ok := reqBody.([]byte)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req, err := changeBytesModelId(req, modelId)
		if err != nil {
			return nil, nil, err
		}
		return ad.CreateImageByBytes(ctx, req)
	case RelayImageEdit:
		req, ok := reqBody.(*adapterApi.EditImageRequest)
		if !ok {
			return nil, nil, errors.New("invalid request body")

		}
		req.Model = modelId
		return ad.CreateImageEdit(ctx, req)
	case RelayImageVariations:
		req, ok := reqBody.(*adapterApi.CreateImageVariationRequest)
		if !ok {
			return nil, nil, errors.New("invalid request body")
		}
		req.Model = modelId
		return ad.ImageVariations(ctx, req)
	}
	return nil, nil, errors.New("invalid relay type")
}

func (s *oaiService) ChatCompletions2Archive(ctx context.Context, req *adapterApi.ChatRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, reqModelId)
		if err != nil {
			s.Logger.Warn("获取provider失败", zap.String("modelId", reqModelId), zap.Error(err))
			continue
		}
		adapterX, err := NewOAIAdapter(conf)
		//newProvider, err := NewOAIProvider(conf)
		if err != nil {
			s.Logger.Warn("创建provider失败", zap.String("modelId", reqModelId), zap.Error(err))
			continue
		}
		req.Model = conf.ModelKey
		resp, respHeader, err := adapterX.ChatCompletions(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			//s.GoLogReq(ctx, conf.ModelKey, 1)
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
	//s.GoLogReq(ctx, reqModelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) ChatCompletionsByBytes2Archive(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, modelId)
		if err != nil {
			continue
		}

		adapterX, err := NewOAIAdapter(conf)
		if err != nil {
			s.Logger.Warn("创建provider失败", zap.String("modelId", modelId), zap.Error(err))
			continue
		}
		req, err = changeBytesModelId(req, conf.ModelId)
		resp, respHeader, err := adapterX.ChatCompletionsByBytes(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			// s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	//s.GoLogReq(ctx, modelId, 2)
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

func (s *oaiService) Completions2Archive(ctx context.Context, req *adapterApi.CompletionsRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, reqModelId)
		if err != nil {
			continue
		}
		adapterX, err := NewOAIAdapter(conf)
		if err != nil {
			s.Logger.Warn("创建provider失败", zap.String("modelId", reqModelId), zap.Error(err))
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := adapterX.Completions(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			// s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	//s.GoLogReq(ctx, reqModelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) CompletionsByBytes2Archive(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, modelId)
		if err != nil {
			continue
		}
		adapterX, err := NewOAIAdapter(conf)
		if err != nil {
			s.Logger.Warn("创建provider失败", zap.String("modelId", modelId), zap.Error(err))
			continue
		}
		req, err = changeBytesModelId(req, conf.ModelId)
		resp, respHeader, err := adapterX.CompletionsByBytes(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			// s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	//s.GoLogReq(ctx, modelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) Embeddings2Archive(ctx context.Context, req *adapterApi.EmbeddingRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, reqModelId)
		if err != nil {
			continue
		}
		adapterX, err := NewOAIAdapter(conf)
		if err != nil {
			s.Logger.Warn("创建provider失败", zap.String("modelId", reqModelId), zap.Error(err))
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := adapterX.Embeddings(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			// s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	//s.GoLogReq(ctx, reqModelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) EmbeddingsByBytes2Archive(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, modelId)
		if err != nil {
			continue
		}
		adapterX, err := NewOAIAdapter(conf)
		if err != nil {
			s.Logger.Warn("创建provider失败", zap.String("modelId", modelId), zap.Error(err))
			continue
		}
		req, err = changeBytesModelId(req, conf.ModelId)
		resp, respHeader, err := adapterX.EmbeddingsByBytes(ctx, req)
		if err == nil {
			// s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	//s.GoLogReq(ctx, modelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) CreateSpeech2Archive(ctx context.Context, req *adapterApi.SpeechRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, reqModelId)
		if err != nil {
			continue
		}
		adapterX, err := NewOAIAdapter(conf)
		if err != nil {
			s.Logger.Warn("创建provider失败", zap.String("modelId", reqModelId), zap.Error(err))
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := adapterX.CreateSpeech(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			// s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	//s.GoLogReq(ctx, reqModelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) CreateSpeechByBytes2Archive(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, modelId)
		if err != nil {
			continue
		}
		adapterX, err := NewOAIAdapter(conf)
		if err != nil {
			s.Logger.Warn("创建provider失败", zap.String("modelId", modelId), zap.Error(err))
			continue
		}
		req, err = changeBytesModelId(req, conf.ModelId)
		resp, respHeader, err := adapterX.CreateSpeechByBytes(ctx, req)
		if err == nil {
			// s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	//s.GoLogReq(ctx, modelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) Transcriptions2Archive(ctx context.Context, req *adapterApi.TranscriptionRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, reqModelId)
		if err != nil {
			continue
		}
		adapterX, err := NewOAIAdapter(conf)
		if err != nil {
			s.Logger.Warn("创建provider失败", zap.String("modelId", reqModelId), zap.Error(err))
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := adapterX.Transcriptions(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			// s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	//s.GoLogReq(ctx, reqModelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) Translations2Archive(ctx context.Context, req *adapterApi.TranslationRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, reqModelId)
		if err != nil {
			continue
		}
		adapterX, err := NewOAIAdapter(conf)
		if err != nil {
			s.Logger.Warn("创建provider失败", zap.String("modelId", reqModelId), zap.Error(err))
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := adapterX.Translations(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			// s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	//s.GoLogReq(ctx, reqModelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) CreateImage2Archive(ctx context.Context, req *adapterApi.CreateImageRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, reqModelId)
		if err != nil {
			continue
		}
		adapterX, err := NewOAIAdapter(conf)
		if err != nil {
			s.Logger.Warn("创建provider失败", zap.String("modelId", reqModelId), zap.Error(err))
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := adapterX.CreateImage(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			// s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	//s.GoLogReq(ctx, reqModelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) CreateImageByBytes2Archive(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, modelId)
		if err != nil {
			continue
		}
		adapterX, err := NewOAIAdapter(conf)
		if err != nil {
			s.Logger.Warn("创建provider失败", zap.String("modelId", modelId), zap.Error(err))
			continue
		}
		req, err = changeBytesModelId(req, conf.ModelId)
		resp, respHeader, err := adapterX.CreateImageByBytes(ctx, req)
		if err == nil {
			// s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		// 标记模型不可用
		s.FailCb(ctx, conf.ModelRecordId)
	}
	//s.GoLogReq(ctx, modelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) CreateImageEdit2Archive(ctx context.Context, req *adapterApi.EditImageRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, reqModelId)
		if err != nil {
			continue
		}
		adapterX, err := NewOAIAdapter(conf)
		if err != nil {
			s.Logger.Warn("创建provider失败", zap.String("modelId", reqModelId), zap.Error(err))
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := adapterX.CreateImageEdit(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			// s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		s.FailCb(ctx, conf.ModelRecordId)
	}
	//s.GoLogReq(ctx, reqModelId, 2)
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) ImageVariations2Archive(ctx context.Context, req *adapterApi.CreateImageVariationRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextChannel(ctx, reqModelId)
		if err != nil {
			continue
		}
		adapterX, err := NewOAIAdapter(conf)
		if err != nil {
			s.Logger.Warn("创建provider失败", zap.String("modelId", reqModelId), zap.Error(err))
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := adapterX.ImageVariations(ctx, req)
		if err == nil {
			s.SuccessCb(ctx, conf.ModelRecordId)
			// s.GoLogReq(ctx, conf.ModelId, 1)
			return resp, respHeader, nil
		}
		s.FailCb(ctx, conf.ModelRecordId)
	}
	//s.GoLogReq(ctx, reqModelId, 2)
	return nil, nil, errors.New("no service available")
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
