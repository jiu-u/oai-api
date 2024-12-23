// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire_load

import (
	"github.com/google/wire"
	"github.com/jiu-u/oai-api/internal/handler"
	"github.com/jiu-u/oai-api/internal/repository"
	"github.com/jiu-u/oai-api/internal/server"
	"github.com/jiu-u/oai-api/internal/service"
	"github.com/jiu-u/oai-api/pkg/config"
	"github.com/jiu-u/oai-api/pkg/log"
	"github.com/jiu-u/oai-api/pkg/sid"
)

// Injectors from wire.go:

func NewWire(cfg *config.Config, logger *log.Logger) (*server.DataLoadTask, func(), error) {
	sidSid := sid.NewSid()
	db := repository.NewDB(cfg)
	repositoryRepository := repository.NewRepository(logger, db)
	transaction := repository.NewTransaction(repositoryRepository)
	serviceService := service.NewService(sidSid, transaction, logger)
	providerRepo := repository.NewProviderRepo(repositoryRepository)
	modelRepo := repository.NewModelRepo(repositoryRepository)
	providerService := service.NewProviderService(serviceService, providerRepo, modelRepo)
	dataLoadTask := server.NewDataLoad(providerService, cfg, logger)
	return dataLoadTask, func() {
	}, nil
}

// wire.go:

var repositorySet = wire.NewSet(repository.NewDB, repository.NewRepository, repository.NewTransaction, repository.NewModelRepo, repository.NewProviderRepo)

var serviceSet = wire.NewSet(service.NewService, service.NewOaiService, service.NewProviderService, service.NewLoadBalanceService)

var handlerSet = wire.NewSet(handler.NewHandler, handler.NewOAIHandler)

var serverSet = wire.NewSet(server.NewDataLoad)
