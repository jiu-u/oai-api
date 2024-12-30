package server

import (
	"context"
	"github.com/jiu-u/oai-api/internal/repository"
	"github.com/jiu-u/oai-api/internal/service"
	adapterV1 "github.com/jiu-u/oai-api/pkg/adapter/api/v1"
	"github.com/jiu-u/oai-api/pkg/config"
	"github.com/jiu-u/oai-api/pkg/log"
	"github.com/lithammer/shortuuid/v4"
	"go.uber.org/zap"
	"io"
	"time"
)

type CheckModelServer struct {
	channelRepo      repository.ChannelRepository
	channelModelRepo repository.ChannelModelRepository
	lbSvc            service.LoadBalanceServiceBeta
	cfg              *config.Config
	logger           *log.Logger
}

func NewCheckModelServer(
	cfg *config.Config,
	lbSvc service.LoadBalanceServiceBeta,
	channelRepo repository.ChannelRepository,
	channelModelRepo repository.ChannelModelRepository,
	logger *log.Logger,
) *CheckModelServer {
	return &CheckModelServer{
		lbSvc:            lbSvc,
		channelRepo:      channelRepo,
		channelModelRepo: channelModelRepo,
		cfg:              cfg,
		logger:           logger,
	}
}

func (c *CheckModelServer) Start(ctx context.Context) error {
	return c.CheckModelChatStatus(ctx)
}

func (c *CheckModelServer) Stop(ctx context.Context) error {
	return nil
}

func (c *CheckModelServer) CheckModelChatStatus(ctx context.Context) error {
	checklist := c.cfg.ChatCompletionCheck
	task := make(chan string, 1)
	go func() {
		for _ = range task {
			uid := shortuuid.New()
			ctx = c.logger.WithValue(context.Background(), zap.String("traceId", uid), zap.String("type", "check_cron"))
			err := c.lbSvc.RecoverChannelModels(ctx)
			if err != nil {
				c.logger.WithContext(ctx).Error("定时检查|chat|模型恢复失败", zap.Error(err))
			}
			c.logger.WithContext(ctx).Info("一轮定时检查开始")
			for _, modelId := range checklist {
				modelIds := []string{modelId}
				if _, ok := c.cfg.ModelMapping[modelId]; ok {
					modelIds = append(modelIds, c.cfg.ModelMapping[modelId]...)
				}
				err := c.CheckModel(ctx, modelIds)
				if err != nil {
					c.logger.WithContext(ctx).Warn("定时检查|chat|"+modelId+"|失败", zap.Error(err))
					continue
				}
				c.logger.WithContext(ctx).Info("定时检查|chat|" + modelId + "|完成")
			}
			c.logger.WithContext(ctx).Info("一轮定时检查完成")
			time.Sleep(60 * time.Minute)
			task <- "ok"
		}
	}()
	task <- "ok"
	return nil
}

func (c *CheckModelServer) CheckModel(ctx context.Context, modelIds []string) error {
	list, err := c.channelModelRepo.FindCheckChannelModels(ctx, modelIds)
	if err != nil {
		return err
	}
	logger := c.logger.WithContext(ctx)
	for _, item := range list {
		time.Sleep(5 * time.Second)
		ctx = context.WithValue(ctx, "modelId", item.ModelKey)
		zapLogger := logger.With(
			zap.Uint64("channelId", item.ChannelId),
			zap.String("modelKey", item.ModelKey),
			zap.Uint64("modelRecordId", item.Id),
		)
		channel, err := c.channelRepo.FindChannelById(ctx, item.ChannelId)
		if err != nil {
			zapLogger.Error("定时检查|chat|数据库根据id获取channel失败", zap.Error(err))
			continue
		}
		zapLogger = zapLogger.With(
			zap.String("channelName", channel.Name),
			zap.String("providerApiKey", channel.APIKey),
		)
		conf := &service.ChannelModelConf{
			ChannelId:       item.ChannelId,
			ChannelName:     channel.Name,
			ChannelType:     channel.Type,
			ChannelKey:      channel.APIKey,
			ChannelEndPoint: channel.EndPoint,
			ModelRecordId:   item.Id,
			ModelKey:        item.ModelKey,
			Weight:          item.Weight,
		}
		newProvider, err := service.NewOAIProvider(conf)
		if err != nil {
			zapLogger.Warn("定时检查|chat|创建provider失败", zap.Error(err))
			continue
		}
		body, _, err := newProvider.ChatCompletions(ctx, &adapterV1.ChatCompletionRequest{
			Model: item.ModelKey,
			Messages: []adapterV1.Message{
				{
					Role:    "user",
					Content: []byte(`"hello,测试!"`),
				},
			},
			Stream:    true,
			MaxTokens: 10,
		})
		if err != nil {
			if body != nil {
				bodyDetail, err := io.ReadAll(body)
				if err != nil {
					zapLogger.Warn("定时检查|chat|对话请求失败|读取body失败", zap.Error(err))
				} else {
					zapLogger.Warn("定时检查|chat|对话请求失败", zap.Error(err), zap.String("detail", string(bodyDetail)))
				}
			} else {
				zapLogger.Warn("定时检查|chat|对话请求失败", zap.Error(err), zap.String("detail", err.Error()))
			}
			// 标记模型不可用
			err = c.lbSvc.FailCb(ctx, item.Id)
			if err != nil {
				zapLogger.Warn("定时检查|chat|更新模型状态失败", zap.Error(err))
			}
			continue
		}
		err = c.lbSvc.SuccessCb(ctx, item.Id)
		if err != nil {
			zapLogger.Warn("定时检查|chat|更新模型状态失败", zap.Error(err))
		}
		zapLogger.Info("定时检查|chat|对话请求成功")
	}
	return nil
}
