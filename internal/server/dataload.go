package server

import (
	"context"
	"fmt"
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/service"
	"github.com/jiu-u/oai-api/pkg/config"
	"os"
	"strings"
)

type DataLoadTask struct {
	Cfg *config.Config
	svc service.ProviderService
}

func NewDataLoad(svc service.ProviderService, cfg *config.Config) *DataLoadTask {
	return &DataLoadTask{
		svc: svc,
		Cfg: cfg,
	}
}

func (s *DataLoadTask) Start(ctx context.Context) error {
	return s.run(ctx)
}

func (s *DataLoadTask) Stop(ctx context.Context) error {
	return nil
}

func (s *DataLoadTask) run(ctx context.Context) error {
	providers := s.Cfg.Providers
	total := 0
	succ := 0
	repeat := 0
	for _, provider := range providers {
		for _, key := range provider.APIKeys {
			total++
			conf := &service.ProviderConf{
				ProviderName:   provider.Name,
				ProviderType:   provider.Type,
				EndPoint:       provider.EndPoint,
				APIKey:         key,
				ProviderModels: provider.Models,
			}
			instance, err := service.NewProvider(conf)
			if err != nil {
				continue
			}
			models := make([]string, 0)
			modelSet := make(map[string]bool)
			getModels, err := instance.Models(ctx)
			if err != nil {
				fmt.Println("获取模型失败", err)
				continue
			}
			for _, model := range getModels {
				if _, ok := modelSet[model]; !ok {
					models = append(models, model)
					modelSet[model] = true
				}
			}
			for _, model := range provider.Models {
				if _, ok := modelSet[model]; !ok {
					models = append(models, model)
					modelSet[model] = true
				}
			}
			req := apiV1.CreateProviderRequest{
				Name:     conf.ProviderName,
				Type:     conf.ProviderType,
				EndPoint: conf.EndPoint,
				APIKey:   key,
				Weight:   provider.Weight,
				Models:   models,
			}
			_, err = s.svc.CreateProvider(ctx, &req)
			if err != nil {
				if strings.Contains(err.Error(), "provider already exists") {
					repeat++
					succ++
					continue
				}
				fmt.Println("创建provider失败", err)
				continue
			}
			succ++
		}
	}
	fmt.Println("总共创建：", total, "个provider")
	fmt.Println("成功创建：", succ, "个provider")
	fmt.Println("重复创建：", repeat, "个provider")
	os.Exit(0)
	return nil
}
