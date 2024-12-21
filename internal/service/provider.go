package service

import (
	"context"
	"errors"
	"fmt"
	v1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/model"
	"github.com/jiu-u/oai-api/internal/repository"
	"github.com/jiu-u/oai-api/pkg/encrypte"
)

type ProviderService interface {
	CreateProvider(ctx context.Context, provider *v1.CreateProviderRequest) (uint64, error)
}

func NewProviderService(srv *Service, repo repository.ProviderRepo, modelRepo repository.ModelRepo) ProviderService {
	return &providerService{
		Service:   srv,
		repo:      repo,
		modelRepo: modelRepo,
	}
}

type providerService struct {
	*Service
	repo      repository.ProviderRepo
	modelRepo repository.ModelRepo
}

func (s *providerService) CreateProvider(ctx context.Context, req *v1.CreateProviderRequest) (uint64, error) {
	id := s.sid.GenUint64()
	format := "%s@%s@%s"
	hashId := encrypte.Sha256Encode(fmt.Sprintf(format, req.Type, req.EndPoint, req.APIKey))
	provider := &model.Provider{
		Name:     req.Name,
		Type:     req.Type,
		EndPoint: req.EndPoint,
		APIKey:   req.APIKey,
		HashId:   hashId,
	}
	provider.Id = id
	err := s.tm.Transaction(ctx, func(ctx context.Context) error {
		exist, _ := s.repo.ExistsHashId(ctx, provider.HashId)
		if exist {
			return errors.New("provider already exists")
		}
		err := s.repo.InsertOne(ctx, provider)
		if err != nil {
			return err
		}

		for _, modelId := range req.Models {
			fmt.Println("modelId--->", modelId)
			newId := s.sid.GenUint64()
			newModel := &model.Model{
				ProviderId: id,
				ModelKey:   modelId,
				Weight:     req.Weight,
			}
			newModel.Id = newId
			_ = s.modelRepo.Insert(ctx, newModel)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}
