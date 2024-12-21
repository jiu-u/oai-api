package server

import (
	"context"
	"github.com/jiu-u/oai-api/internal/model"
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
	providerRepo repository.ProviderRepo
	modelRepo    repository.ModelRepo
	cfg          *config.Config
	logger       *log.Logger
}

func NewCheckModelServer(
	modelRepo repository.ModelRepo,
	cfg *config.Config,
	providerRepo repository.ProviderRepo,
	logger *log.Logger,
) *CheckModelServer {
	return &CheckModelServer{
		modelRepo:    modelRepo,
		cfg:          cfg,
		providerRepo: providerRepo,
		logger:       logger,
	}
}

func (c *CheckModelServer) Start(ctx context.Context) error {
	return c.CheckModelChatStatus(ctx)
}

func (c *CheckModelServer) Stop(ctx context.Context) error {
	return nil
}

func (c *CheckModelServer) CheckModelChatStatus(ctx context.Context) error {
	checklist := c.cfg.ChatCOmpletionCheck
	task := make(chan string, 1)
	go func() {
		for _ = range task {
			uid := shortuuid.New()
			ctx = c.logger.WithValue(context.Background(), zap.String("traceId", uid), zap.String("type", "check_cron"))
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
	list, err := c.modelRepo.FindCheckModels(ctx, modelIds)
	if err != nil {
		return err
	}
	logger := c.logger.WithContext(ctx)
	for _, item := range list {
		time.Sleep(5 * time.Second)
		ctx = context.WithValue(ctx, "modelId", item.Model)
		zapLogger := logger.With(
			zap.Uint64("providerId", item.ProviderId),
			zap.String("model", item.Model),
			zap.Uint64("modelId", item.Id),
		)
		//ctx = c.logger.WithValue(
		//	ctx,
		//	zap.Uint64("providerId", item.ProviderId),
		//	zap.String("model", item.Model),
		//	zap.Uint64("modelId", item.Id),
		//)
		providerConf, err := c.providerRepo.FindOne(ctx, item.ProviderId)
		if err != nil {
			zapLogger.Warn("定时检查|chat|根据id获取provider失败", zap.Error(err))
			continue
		}
		zapLogger = zapLogger.With(
			zap.String("providerName", providerConf.Name),
			zap.String("providerApiKey", providerConf.APIKey),
		)
		//ctx := c.logger.WithValue(ctx,
		//	zap.String("providerName", providerConf.Name),
		//	zap.String("providerApiKey", providerConf.APIKey),
		//)
		conf := &service.ProviderConf{
			ProviderName: providerConf.Name,
			ProviderType: providerConf.Type,
			EndPoint:     providerConf.EndPoint,
			APIKey:       providerConf.APIKey,
		}
		newProvider, err := service.NewProvider(conf)
		if err != nil {
			zapLogger.Warn("定时检查|chat|创建provider失败", zap.Error(err))
			continue
		}
		body, _, err := newProvider.ChatCompletions(ctx, &adapterV1.ChatCompletionRequest{
			Model: item.Model,
			Messages: []adapterV1.Message{
				{
					Role:    "user",
					Content: "hello,测试!",
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
					continue
				}
				zapLogger.Warn("定时检查|chat|对话请求失败", zap.Error(err), zap.String("detail", string(bodyDetail)))
			}
			// 标记模型不可用
			zapLogger.Warn("定时检查|chat|对话请求失败", zap.Error(err), zap.String("detail", err.Error()))
			err = c.modelRepo.UpdateStatus(ctx, item.Id, model.ChatStatus, 0)
			if err != nil {
				zapLogger.Warn("定时检查|chat|更新模型状态失败", zap.Error(err))
			}
			continue
		}
		err = c.modelRepo.UpdateStatus(ctx, item.Id, model.ChatStatus, 1)
		if err != nil {
			zapLogger.Warn("定时检查|chat|更新模型状态失败", zap.Error(err))
		}
		zapLogger.Info("定时检查|chat|对话请求成功")
	}
	return nil
}
