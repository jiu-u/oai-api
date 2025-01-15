//go:build wireinject
// +build wireinject

package wire_load

import (
	"github.com/google/wire"
	"github.com/jiu-u/oai-api/internal/handler"
	"github.com/jiu-u/oai-api/internal/repository"
	"github.com/jiu-u/oai-api/internal/server"
	"github.com/jiu-u/oai-api/internal/service"
	"github.com/jiu-u/oai-api/pkg/cache"
	"github.com/jiu-u/oai-api/pkg/config"
	"github.com/jiu-u/oai-api/pkg/jwt"
	"github.com/jiu-u/oai-api/pkg/log"
	"github.com/jiu-u/oai-api/pkg/sid"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewChannelModelRepository,
	repository.NewChannelRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewOaiService,
	service.NewChannelService,
	service.NewLoadBalanceServiceBeta,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewOAIHandler,
)

var serverSet = wire.NewSet(
	server.NewDataLoad,
)

func NewWire(cfg *config.Config, logger *log.Logger) (*server.DataLoadTask, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		serverSet,
		sid.NewSid,
		jwt.NewJwt,
		cache.New,
	))
}
