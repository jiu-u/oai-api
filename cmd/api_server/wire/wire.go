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
	repository.NewUserRepository,
	repository.NewApiKeyRepository,
	repository.NewRequestLogRepository,
	repository.NewChannelRepository,
	repository.NewChannelModelRepository,
	repository.NewSystemRepository,
	repository.NewUserAuthProviderRepository,
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
	service.NewSystemConfigService,
	service.NewEmailService,
	service.NewVerificationService,
	service.NewModelCheckService,
	oauth2.NewLinuxDoAuthService,
	oauth2.NewGithubAuthService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewOAIHandler,
	handler.NewApiKeyHandler,
	handler.NewAuthHandler,
	handler.NewRequestLogHandler,
	handler.NewUserHandler,
	handler.NewSystemConfigHandler,
	handler.NewVerificationHandler,
	handler.NewChannelHandler,
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
		//app.WithServer(httpServer),
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
		jwt.NewJwt,
		cache.New,
		newApp,
	))
}
