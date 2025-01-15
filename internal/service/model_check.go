package service

import (
	"context"
	"fmt"
	adapterApi "github.com/jiu-u/oai-adapter/api"
	"github.com/jiu-u/oai-api/internal/dto"
	"github.com/jiu-u/oai-api/internal/model"
	"github.com/jiu-u/oai-api/internal/repository"
	"go.uber.org/zap"
	"io"
	"strconv"
	"time"
)

type ModelCheckService interface {
	CheckModel(ctx context.Context, conf *dto.ChannelModelConf) (*dto.ModelCheckResult, error)
	CheckModel2(ctx context.Context, channelId uint64, modelId string) (*dto.ModelCheckResult, error)
}

func NewModelCheckService(
	s *Service,
	channelRepo repository.ChannelRepository,
	channelModelRepo repository.ChannelModelRepository,
	lbSvc LoadBalanceServiceBeta,
) ModelCheckService {
	return &modelCheckService{
		Service:          s,
		channelRepo:      channelRepo,
		channelModelRepo: channelModelRepo,
		lbSvc:            lbSvc,
	}
}

type modelCheckService struct {
	*Service
	lbSvc            LoadBalanceServiceBeta
	channelRepo      repository.ChannelRepository
	channelModelRepo repository.ChannelModelRepository
}

func (s *modelCheckService) CheckModel2(ctx context.Context, channelId uint64, modelId string) (*dto.ModelCheckResult, error) {
	channelX, err := s.channelRepo.FindChannelById(ctx, channelId)
	if err != nil {
		return nil, err
	}
	modelX, err := s.channelModelRepo.ExistsChannelModel(ctx, &model.ChannelModel{
		ChannelId: channelId,
		ModelKey:  modelId,
	})
	if err != nil {
		return nil, err
	}
	conf := &dto.ChannelModelConf{
		ChannelId:       channelId,
		ChannelName:     channelX.Name,
		ChannelType:     channelX.Type,
		ChannelKey:      channelX.APIKey,
		ChannelEndPoint: channelX.EndPoint,
		ModelRecordId:   modelX.Id,
		ModelKey:        modelId,
		ModelId:         modelId,
		Weight:          10,
	}
	return s.CheckModel(ctx, conf)
}

func (s *modelCheckService) CheckModel(ctx context.Context, conf *dto.ChannelModelConf) (*dto.ModelCheckResult, error) {
	startTime := time.Now()
	adapterX, err := NewOAIAdapter(conf)
	if err != nil {
		return nil, fmt.Errorf("创建provider失败: %s", err.Error())
	}
	body, _, err := adapterX.ChatCompletions(ctx, &adapterApi.ChatRequest{
		Model: conf.ModelKey,
		Messages: []adapterApi.Message{
			{
				Role:    "user",
				Content: []byte(`"hello,测试!"`),
			},
		},
		Stream:    true,
		MaxTokens: 10,
	})
	if err != nil {
		go func() {
			err := s.lbSvc.FailCb(ctx, conf.ModelRecordId)
			if err != nil {
				s.Logger.Warn("failCb失败", zap.Error(err))
			}
		}()
		if body != nil {
			bodyDetail, err2 := io.ReadAll(body)
			if err2 == nil {
				err = fmt.Errorf("请求失败: %s", string(bodyDetail))
				return nil, err
			}
		}
		return nil, fmt.Errorf("请求失败: %s", err.Error())

	}
	defer body.Close()
	duration := time.Since(startTime)
	go func() {
		err := s.lbSvc.SuccessCb(ctx, conf.ModelRecordId)
		if err != nil {
			s.Logger.Warn("successCb失败", zap.Error(err))
		}
	}()
	return &dto.ModelCheckResult{
		ChannelId:          strconv.FormatUint(conf.ChannelId, 10),
		ModelName:          conf.ModelKey,
		ConnectionDuration: int64(duration.Milliseconds()),
		TotalDuration:      int64(duration.Milliseconds()),
		Status:             1,
	}, nil
}
