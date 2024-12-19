package service

import (
	"context"
	"errors"
	v1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/model"
	"github.com/jiu-u/oai-api/internal/repository"
)

type ProviderService interface {
	CreateProvider(ctx context.Context, provider *v1.CreateProviderRequest) (uint64, error)
}

func NewProviderService(srv *Service, repo repository.ProviderRepo) ProviderService {
	return &providerService{
		Service: srv,
		repo:    repo,
	}
}

type providerService struct {
	*Service
	repo repository.ProviderRepo
}

func (s *providerService) CreateProvider(ctx context.Context, req *v1.CreateProviderRequest) (uint64, error) {
	id := s.sid.GenUint64()
	provider := &model.Provider{
		Name:     req.Name,
		Type:     req.Type,
		EndPoint: req.EndPoint,
		APIKey:   req.APIKey,
	}
	provider.Id = id
	err := s.tm.Transaction(ctx, func(ctx context.Context) error {
		exist, _ := s.repo.ExistsHashId(ctx, provider.HashId)
		if exist {
			return errors.New("provider already exists")
		}
		return s.repo.InsertOne(ctx, provider)
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}
