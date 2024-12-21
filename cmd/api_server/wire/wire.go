//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/jiu-u/oai-api/internal/handler"
	"github.com/jiu-u/oai-api/internal/repository"
	"github.com/jiu-u/oai-api/internal/server"
	"github.com/jiu-u/oai-api/internal/service"
	"github.com/jiu-u/oai-api/pkg/app"
	"github.com/jiu-u/oai-api/pkg/config"
	"github.com/jiu-u/oai-api/pkg/log"
	"github.com/jiu-u/oai-api/pkg/server/http"
	"github.com/jiu-u/oai-api/pkg/sid"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewModelRepo,
	repository.NewProviderRepo,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewOaiService,
	service.NewProviderService,
	service.NewLoadBalanceService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewOAIHandler,
)

var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewCheckModelServer,
)

// build App
func newApp(
	httpServer *http.Server,
	checkServer *server.CheckModelServer,
	// job *server.Job,
	// task *server.Task,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, checkServer),
		app.WithName("demo-server"),
	)
}

func NewWire(cfg *config.Config, logger *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		serverSet,
		sid.NewSid,
		newApp,
	))
}
