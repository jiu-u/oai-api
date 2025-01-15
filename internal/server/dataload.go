package server

import (
	"context"
	"github.com/jiu-u/oai-api/internal/service"
	"github.com/jiu-u/oai-api/pkg/config"
	"github.com/jiu-u/oai-api/pkg/log"
)

type DataLoadTask struct {
	Cfg    *config.Config
	svc    service.ChannelService
	logger *log.Logger
}

func NewDataLoad(svc service.ChannelService, cfg *config.Config, logger *log.Logger) *DataLoadTask {
	return &DataLoadTask{
		svc:    svc,
		Cfg:    cfg,
		logger: logger,
	}
}

func (s *DataLoadTask) Start(ctx context.Context) error {
	return s.run(ctx)
}

func (s *DataLoadTask) Stop(ctx context.Context) error {
	return nil
}

func (s *DataLoadTask) run(ctx context.Context) error {
	return nil
	//providers := s.Cfg.Providers
	//total := 0
	//succ := 0
	////repeat := 0
	//for _, provider := range providers {
	//	for _, key := range provider.APIKeys {
	//		total++
	//		conf := &dto.ChannelModelConf{
	//			ChannelName:     provider.Name,
	//			ChannelType:     provider.Type,
	//			ChannelKey:      key,
	//			ChannelEndPoint: provider.EndPoint,
	//		}
	//		instance, err := service.NewOAIAdapter(conf)
	//		if err != nil {
	//			continue
	//		}
	//		models := make([]string, 0)
	//		modelSet := make(map[string]bool)
	//		getModels, err := instance.Models(ctx)
	//		if err != nil {
	//			s.logger.Error("获取模型失败", zap.Error(err), zap.Any("channel", conf))
	//			continue
	//		}
	//		for _, model := range getModels {
	//			if _, ok := modelSet[model]; !ok {
	//				models = append(models, model)
	//				modelSet[model] = true
	//			}
	//		}
	//		for _, model := range provider.Models {
	//			if _, ok := modelSet[model]; !ok {
	//				models = append(models, model)
	//				modelSet[model] = true
	//			}
	//		}
	//		req := apiV1.CreateChannelRequest{
	//			Name:     conf.ChannelName,
	//			Type:     conf.ChannelType,
	//			EndPoint: conf.ChannelEndPoint,
	//			APIKey:   key,
	//			Weight:   provider.Weight,
	//			Models:   models,
	//		}
	//		_, err = s.svc.CreateChannel(ctx, &req)
	//		if err != nil {
	//			s.logger.Error("创建provider失败", zap.Error(err))
	//			continue
	//		}
	//		succ++
	//	}
	//}
	//s.logger.Info("总共创建：" + strconv.Itoa(total) + "个provider")
	//s.logger.Info("成功创建：" + strconv.Itoa(succ) + "个provider")
	////s.logger.Info("重复数据：" + strconv.Itoa(repeat) + "个provider")
	//s.logger.Info("数据加载完成")
	//return nil
}
