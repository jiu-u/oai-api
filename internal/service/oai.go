package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/model"
	"github.com/jiu-u/oai-api/internal/repository"
	adapterV1 "github.com/jiu-u/oai-api/pkg/adapter/api/v1"
	"github.com/jiu-u/oai-api/pkg/adapter/provider"
	"github.com/jiu-u/oai-api/pkg/array"
	"io"
	"net/http"
)

type OaiService interface {
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

func NewProvider(config *ProviderConf) (provider.Provider, error) {
	conf := provider.Config{
		Type:     config.ProviderType,
		EndPoint: config.EndPoint,
		APIKey:   config.APIKey,
	}
	switch conf.Type {
	case "openai":
		return provider.NewOpenAIProvider(conf), nil
	case "oaiNoFetchModel":
		return provider.NewOaiNoFetchModelProvider(conf, config.ProviderModels), nil
	case "siliconflow":
		return provider.NewSiliconFlowProvider(conf), nil
	default:
		return nil, errors.New("invalid provider type")
	}
}

func NewOaiService(
	svc *Service,
	load LoadBalanceService,
	modelRepo repository.ModelRepo,
) OaiService {
	return &oaiService{
		Service:   svc,
		load:      load,
		modelRepo: modelRepo,
		N:         2,
	}
}

type oaiService struct {
	*Service
	load      LoadBalanceService
	modelRepo repository.ModelRepo
	N         int
}

func (s *oaiService) ChatCompletions(ctx context.Context, req *adapterV1.ChatCompletionRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextProvider(ctx, reqModelId, model.ChatStatus)
		if err != nil {
			fmt.Println("获取provider失败", err)
			continue
		}
		newProvider, err := NewProvider(conf)
		if err != nil {
			fmt.Println("创建provider失败", err)
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := newProvider.ChatCompletions(ctx, req)
		if err == nil {
			return resp, respHeader, nil
		}
		fmt.Println("获取response失败", err)
		// 标记模型不可用
		err = s.modelRepo.UpdateStatus(ctx, conf.ModelUID, model.ChatStatus, 0)
		fmt.Println("更新状态失败", err)
	}
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) ChatCompletionsByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	for _ = range s.N {
		conf, err := s.load.NextProvider(ctx, modelId, model.ChatStatus)
		if err != nil {
			continue
		}
		newProvider, err := NewProvider(conf)
		if err != nil {
			continue
		}
		req, err = changeBytesModelId(req, conf.ModelId)
		resp, respHeader, err := newProvider.ChatCompletionsByBytes(ctx, req)
		if err == nil {
			return resp, respHeader, nil
		}
		// 标记模型不可用
		err = s.modelRepo.UpdateStatus(ctx, conf.ModelUID, model.ChatStatus, 0)
		fmt.Println("更新状态失败", err)
	}
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) Models(ctx context.Context) (*apiV1.ModelResponse, error) {
	modelIds, err := s.modelRepo.GetAllModelIds(ctx)
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

func (s *oaiService) Completions(ctx context.Context, req *adapterV1.CompletionsRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextProvider(ctx, reqModelId, model.ChatStatus)
		if err != nil {
			continue
		}
		newProvider, err := NewProvider(conf)
		if err != nil {
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := newProvider.Completions(ctx, req)
		if err == nil {
			return resp, respHeader, nil
		}
		// 标记模型不可用
		err = s.modelRepo.UpdateStatus(ctx, conf.ModelUID, model.ChatStatus, 0)
		fmt.Println("更新状态失败", err)
	}
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) CompletionsByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	for _ = range s.N {
		conf, err := s.load.NextProvider(ctx, modelId, model.ChatStatus)
		if err != nil {
			continue
		}
		newProvider, err := NewProvider(conf)
		if err != nil {
			continue
		}
		req, err = changeBytesModelId(req, conf.ModelId)
		resp, respHeader, err := newProvider.CompletionsByBytes(ctx, req)
		if err == nil {
			return resp, respHeader, nil
		}
		// 标记模型不可用
		err = s.modelRepo.UpdateStatus(ctx, conf.ModelUID, model.ChatStatus, 0)
		fmt.Println("更新状态失败", err)
	}
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) Embeddings(ctx context.Context, req *adapterV1.EmbeddingRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextProvider(ctx, reqModelId, model.ChatStatus)
		if err != nil {
			continue
		}
		newProvider, err := NewProvider(conf)
		if err != nil {
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := newProvider.Embeddings(ctx, req)
		if err == nil {
			return resp, respHeader, nil
		}
		// 标记模型不可用
		err = s.modelRepo.UpdateStatus(ctx, conf.ModelUID, model.ChatStatus, 0)
		fmt.Println("更新状态失败", err)
	}
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) EmbeddingsByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	for _ = range s.N {
		conf, err := s.load.NextProvider(ctx, modelId, model.ChatStatus)
		if err != nil {
			continue
		}
		newProvider, err := NewProvider(conf)
		if err != nil {
			continue
		}
		req, err = changeBytesModelId(req, conf.ModelId)
		resp, respHeader, err := newProvider.EmbeddingsByBytes(ctx, req)
		if err == nil {
			return resp, respHeader, nil
		}
		// 标记模型不可用
		err = s.modelRepo.UpdateStatus(ctx, conf.ModelUID, model.ChatStatus, 0)
		fmt.Println("更新状态失败", err)
	}
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) CreateSpeech(ctx context.Context, req *adapterV1.SpeechRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextProvider(ctx, reqModelId, model.ChatStatus)
		if err != nil {
			continue
		}
		newProvider, err := NewProvider(conf)
		if err != nil {
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := newProvider.CreateSpeech(ctx, req)
		if err == nil {
			return resp, respHeader, nil
		}
		// 标记模型不可用
		err = s.modelRepo.UpdateStatus(ctx, conf.ModelUID, model.ChatStatus, 0)
		fmt.Println("更新状态失败", err)
	}
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) CreateSpeechByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	for _ = range s.N {
		conf, err := s.load.NextProvider(ctx, modelId, model.ChatStatus)
		if err != nil {
			continue
		}
		newProvider, err := NewProvider(conf)
		if err != nil {
			continue
		}
		req, err = changeBytesModelId(req, conf.ModelId)
		resp, respHeader, err := newProvider.CreateSpeechByBytes(ctx, req)
		if err == nil {
			return resp, respHeader, nil
		}
		// 标记模型不可用
		err = s.modelRepo.UpdateStatus(ctx, conf.ModelUID, model.ChatStatus, 0)
		fmt.Println("更新状态失败", err)
	}
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) Transcriptions(ctx context.Context, req *adapterV1.TranscriptionRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextProvider(ctx, reqModelId, model.ChatStatus)
		if err != nil {
			continue
		}
		newProvider, err := NewProvider(conf)
		if err != nil {
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := newProvider.Transcriptions(ctx, req)
		if err == nil {
			return resp, respHeader, nil
		}
		// 标记模型不可用
		err = s.modelRepo.UpdateStatus(ctx, conf.ModelUID, model.ChatStatus, 0)
		fmt.Println("更新状态失败", err)
	}
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) Translations(ctx context.Context, req *adapterV1.TranslationRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextProvider(ctx, reqModelId, model.ChatStatus)
		if err != nil {
			continue
		}
		newProvider, err := NewProvider(conf)
		if err != nil {
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := newProvider.Translations(ctx, req)
		if err == nil {
			return resp, respHeader, nil
		}
		// 标记模型不可用
		err = s.modelRepo.UpdateStatus(ctx, conf.ModelUID, model.ChatStatus, 0)
		fmt.Println("更新状态失败", err)
	}
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) CreateImage(ctx context.Context, req *adapterV1.CreateImageRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextProvider(ctx, reqModelId, model.ChatStatus)
		if err != nil {
			continue
		}
		newProvider, err := NewProvider(conf)
		if err != nil {
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := newProvider.CreateImage(ctx, req)
		if err == nil {
			return resp, respHeader, nil
		}
		// 标记模型不可用
		err = s.modelRepo.UpdateStatus(ctx, conf.ModelUID, model.ChatStatus, 0)
		fmt.Println("更新状态失败", err)
	}
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) CreateImageByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	for _ = range s.N {
		conf, err := s.load.NextProvider(ctx, modelId, model.ChatStatus)
		if err != nil {
			continue
		}
		newProvider, err := NewProvider(conf)
		if err != nil {
			continue
		}
		req, err = changeBytesModelId(req, conf.ModelId)
		resp, respHeader, err := newProvider.CreateImageByBytes(ctx, req)
		if err == nil {
			return resp, respHeader, nil
		}
		// 标记模型不可用
		err = s.modelRepo.UpdateStatus(ctx, conf.ModelUID, model.ChatStatus, 0)
		fmt.Println("更新状态失败", err)
	}
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) CreateImageEdit(ctx context.Context, req *adapterV1.EditImageRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextProvider(ctx, reqModelId, model.ChatStatus)
		if err != nil {
			continue
		}
		newProvider, err := NewProvider(conf)
		if err != nil {
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := newProvider.CreateImageEdit(ctx, req)
		if err == nil {
			return resp, respHeader, nil
		}
		// 标记模型不可用
		err = s.modelRepo.UpdateStatus(ctx, conf.ModelUID, model.ChatStatus, 0)
		fmt.Println("更新状态失败", err)
	}
	return nil, nil, errors.New("no service available")
}

func (s *oaiService) ImageVariations(ctx context.Context, req *adapterV1.CreateImageVariationRequest) (io.ReadCloser, http.Header, error) {
	reqModelId := req.Model
	for _ = range s.N {
		conf, err := s.load.NextProvider(ctx, reqModelId, model.ChatStatus)
		if err != nil {
			continue
		}
		newProvider, err := NewProvider(conf)
		if err != nil {
			continue
		}
		req.Model = conf.ModelId
		resp, respHeader, err := newProvider.ImageVariations(ctx, req)
		if err == nil {
			return resp, respHeader, nil
		}
		// 标记模型不可用
		err = s.modelRepo.UpdateStatus(ctx, conf.ModelUID, model.ChatStatus, 0)
		fmt.Println("更新状态失败", err)
	}
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
