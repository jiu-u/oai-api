package route

import (
	"github.com/jiu-u/oai-api/internal/handler"
	"github.com/jiu-u/oai-api/internal/route/routes"
	"github.com/jiu-u/oai-api/internal/service"
	"github.com/jiu-u/oai-api/pkg/jwt"
	"github.com/jiu-u/oai-api/pkg/log"
	"github.com/jiu-u/oai-api/pkg/server/http"
)

func SetupRoute(
	s *http.Server,
	verificationHandler *handler.VerificationHandler,
	sysConfigHandler *handler.SystemConfigHandler,
	authHandler *handler.AuthHandler,
	apiKeyHandler *handler.ApiKeyHandler,
	oaiHandler *handler.OAIHandler,
	channelHandler *handler.ChannelHandler,
	requestLogHandler *handler.RequestLogHandler,
	userHandler *handler.UserHandler,
	apiKeySvc service.ApiKeyService,
	logger *log.Logger,
	jwtJWT *jwt.JWT,
) {
	v1Group := s.Group("/v1")
	v1BetaGroup := s.Group("/v1beta")
	// 用户登录、注册、登出
	routes.SetupAuthRoutes(v1Group, authHandler, sysConfigHandler, userHandler, jwtJWT, logger)
	// 验证码发送
	routes.SetupVerificationRoutes(v1Group, verificationHandler)
	// 系统配置
	routes.SetupSystemConfigRoutes(v1Group, sysConfigHandler, jwtJWT, logger)
	// oai
	routes.SetupOaiRoutes(v1Group, v1BetaGroup, apiKeyHandler, oaiHandler, apiKeySvc, logger)
	// channel
	routes.SetupChannelRoutes(v1Group, channelHandler, jwtJWT, logger)
	// request log
	routes.SetupOaiReqLogRoutes(v1Group, requestLogHandler, jwtJWT, logger)
	// api key
	routes.SetupApiKeyRoutes(v1Group, apiKeyHandler, jwtJWT, logger)
}
