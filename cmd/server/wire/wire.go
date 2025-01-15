//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/jiu-u/oai-api/internal/handler"
	"github.com/jiu-u/oai-api/internal/repository"
	"github.com/jiu-u/oai-api/internal/server"
	"github.com/jiu-u/oai-api/internal/service"
	"github.com/jiu-u/oai-api/internal/service/oauth2"
	"github.com/jiu-u/oai-api/pkg/app"
	"github.com/jiu-u/oai-api/pkg/cache"
	"github.com/jiu-u/oai-api/pkg/config"
	"github.com/jiu-u/oai-api/pkg/jwt"
	"github.com/jiu-u/oai-api/pkg/log"
	"github.com/jiu-u/oai-api/pkg/server/http"
	"github.com/jiu-u/oai-api/pkg/sid"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewChannelRepository,
	repository.NewChannelModelRepository,
	repository.NewUserRepository,
	repository.NewApiKeyRepository,
	repository.NewRequestLogRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewOaiService,
	service.NewChannelService,
	service.NewLoadBalanceServiceBeta,
	service.NewRequestLogService,
	service.NewApiKeyService,
	service.NewUserService,
	service.NewAuthService,
	oauth2.NewService,
	oauth2.NewLinuxDoAuthService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewOAIHandler,
	handler.NewApiKeyHandler,
	handler.NewAuthHandler,
	handler.NewRequestLogHandler,
	handler.NewUserHandler,
)

var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewCheckModelServer,
	server.NewMigrate,
	server.NewDataLoad,
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

func newWireApp(app *app.App, migrateJob *server.Migrate, dataLoadJob *server.DataLoadTask) *WireApp {
	return &WireApp{
		App:         app,
		MigrateJob:  migrateJob,
		DataLoadJob: dataLoadJob,
	}
}

type WireApp struct {
	App         *app.App
	MigrateJob  *server.Migrate
	DataLoadJob *server.DataLoadTask
}

//func NewMigrateWire(cfg *config.Config, logger *log.Logger) (*server.Migrate, func(), error) {
//	panic(wire.Build(
//		repositorySet,
//		server.NewMigrate,
//	))
//}
//
//func NewDataLoadWire(cfg *config.Config, logger *log.Logger) (*server.DataLoadTask, func(), error) {
//	panic(wire.Build(
//		repositorySet,
//		service.NewService,
//		service.NewChannelService,
//		sid.NewSid,
//		server.NewDataLoad,
//	))
//}

func NewWire(cfg *config.Config, logger *log.Logger) (*WireApp, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		serverSet,
		sid.NewSid,
		jwt.NewJwt,
		cache.New,
		newApp,
		newWireApp,
	))
}
