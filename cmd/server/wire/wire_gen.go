// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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

// Injectors from wire.go:

func NewWire(cfg *config.Config, logger *log.Logger) (*WireApp, func(), error) {
	jwtJWT := jwt.NewJwt(cfg)
	sidSid := sid.NewSid()
	db := repository.NewDB(cfg)
	repositoryRepository := repository.NewRepository(logger, db)
	transaction := repository.NewTransaction(repositoryRepository)
	cacheCache := cache.New()
	serviceService := service.NewService(sidSid, transaction, logger, jwtJWT, cacheCache)
	channelRepository := repository.NewChannelRepository(repositoryRepository)
	channelModelRepository := repository.NewChannelModelRepository(repositoryRepository)
	loadBalanceServiceBeta := service.NewLoadBalanceServiceBeta(serviceService, channelRepository, channelModelRepository)
	userRepository := repository.NewUserRepository(repositoryRepository)
	requestLogRepository := repository.NewRequestLogRepository(repositoryRepository)
	apiKeyRepository := repository.NewApiKeyRepository(repositoryRepository)
	requestLogService := service.NewRequestLogService(serviceService, userRepository, requestLogRepository, apiKeyRepository)
	oaiService := service.NewOaiService(serviceService, loadBalanceServiceBeta, requestLogService, channelModelRepository)
	oaiHandler := handler.NewOAIHandler(oaiService)
	handlerHandler := handler.NewHandler(logger)
	systemRepository := repository.NewSystemRepository(repositoryRepository)
	systemConfigService := service.NewSystemConfigService(serviceService, systemRepository)
	linuxDoOauthService := oauth2.NewLinuxDoAuthService(systemRepository)
	gitHubOauthService := oauth2.NewGithubAuthService(systemRepository)
	userAuthProviderRepository := repository.NewUserAuthProviderRepository(repositoryRepository)
	authService := service.NewAuthService(serviceService, userRepository, systemConfigService, linuxDoOauthService, gitHubOauthService, userAuthProviderRepository)
	authHandler := handler.NewAuthHandler(handlerHandler, jwtJWT, authService, systemConfigService)
	apiKeyService := service.NewApiKeyService(serviceService, userRepository, apiKeyRepository)
	apiKeyHandler := handler.NewApiKeyHandler(handlerHandler, apiKeyService)
	userService := service.NewUserService(serviceService, userRepository, apiKeyRepository)
	userHandler := handler.NewUserHandler(handlerHandler, userService)
	requestLogHandler := handler.NewRequestLogHandler(handlerHandler, requestLogService)
	systemConfigHandler := handler.NewSystemConfigHandler(handlerHandler, systemConfigService)
	emailService := service.NewEmailService(systemRepository)
	verificationService := service.NewVerificationService(serviceService, emailService)
	verificationHandler := handler.NewVerificationHandler(handlerHandler, verificationService)
	channelService := service.NewChannelService(serviceService, channelRepository, channelModelRepository, loadBalanceServiceBeta)
	modelCheckService := service.NewModelCheckService(serviceService, channelRepository, channelModelRepository, loadBalanceServiceBeta)
	channelHandler := handler.NewChannelHandler(handlerHandler, channelService, modelCheckService)
	httpServer := server.NewHTTPServer(logger, cfg, jwtJWT, oaiHandler, authHandler, apiKeyService, apiKeyHandler, userHandler, requestLogHandler, systemConfigHandler, verificationHandler, channelHandler)
	checkModelServer := server.NewCheckModelServer(loadBalanceServiceBeta, channelRepository, channelModelRepository, logger, systemConfigService)
	app := newApp(httpServer, checkModelServer)
	migrate := server.NewMigrate(db, logger)
	wireApp := newWireApp(app, migrate)
	return wireApp, func() {
	}, nil
}

// wire.go:

var repositorySet = wire.NewSet(repository.NewDB, repository.NewRepository, repository.NewTransaction, repository.NewUserRepository, repository.NewApiKeyRepository, repository.NewRequestLogRepository, repository.NewChannelRepository, repository.NewChannelModelRepository, repository.NewSystemRepository, repository.NewUserAuthProviderRepository)

var serviceSet = wire.NewSet(service.NewService, service.NewOaiService, service.NewChannelService, service.NewLoadBalanceServiceBeta, service.NewRequestLogService, service.NewApiKeyService, service.NewUserService, service.NewAuthService, service.NewSystemConfigService, service.NewEmailService, service.NewVerificationService, service.NewModelCheckService, oauth2.NewLinuxDoAuthService, oauth2.NewGithubAuthService)

var handlerSet = wire.NewSet(handler.NewHandler, handler.NewOAIHandler, handler.NewApiKeyHandler, handler.NewAuthHandler, handler.NewRequestLogHandler, handler.NewUserHandler, handler.NewSystemConfigHandler, handler.NewVerificationHandler, handler.NewChannelHandler)

var serverSet = wire.NewSet(server.NewHTTPServer, server.NewCheckModelServer, server.NewMigrate)

// build App
func newApp(
	httpServer *http.Server,
	checkServer *server.CheckModelServer,

) *app.App {
	return app.NewApp(app.WithServer(httpServer, checkServer), app.WithName("demo-server"))
}

func newWireApp(app2 *app.App, migrateJob *server.Migrate) *WireApp {
	return &WireApp{
		App:        app2,
		MigrateJob: migrateJob,
	}
}

type WireApp struct {
	App        *app.App
	MigrateJob *server.Migrate
}
