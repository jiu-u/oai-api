package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jiu-u/oai-api/internal/handler"
	"github.com/jiu-u/oai-api/internal/middleware"
	"github.com/jiu-u/oai-api/internal/route"
	"github.com/jiu-u/oai-api/internal/service"
	"github.com/jiu-u/oai-api/pkg/config"
	"github.com/jiu-u/oai-api/pkg/jwt"
	"github.com/jiu-u/oai-api/pkg/log"
	"github.com/jiu-u/oai-api/pkg/server/http"
)

func NewHTTPServer(
	logger *log.Logger,
	cfg *config.Config,
	jwt2 *jwt.JWT,
	oaiHandler *handler.OAIHandler,
	authHandler *handler.AuthHandler,
	apiKeySvc service.ApiKeyService,
	apiKeyHandler *handler.ApiKeyHandler,
	userHandler *handler.UserHandler,
	requestLogHandler *handler.RequestLogHandler,
	sysConfigHandler *handler.SystemConfigHandler,
	verificationHandler *handler.VerificationHandler,
	channelHandler *handler.ChannelHandler,
) *http.Server {
	//gin.SetMode(gin.DebugMode)
	s := http.NewServer(
		gin.Default(),
		logger,
		http.WithServerHost(cfg.HTTP.Host),
		http.WithServerPort(cfg.HTTP.Port),
	)
	s.Use(
		middleware.TraceMiddleware(logger),
		middleware.CORSMiddleware(),
		middleware.SessionMiddleware(),
	)
	//s.Static("/assets", "./web/dist/assets")
	//s.GET("/", func(ctx *gin.Context) {
	//	ctx.File("./web/dist/index.html")
	//})
	route.SetupRoute(
		s,
		verificationHandler,
		sysConfigHandler,
		authHandler,
		apiKeyHandler,
		oaiHandler,
		channelHandler,
		requestLogHandler,
		userHandler,
		apiKeySvc,
		logger,
		jwt2,
	)

	//apiGroup := s.Group("/api")
	//v1Group := apiGroup.Group("/v1")
	//v1Beta := apiGroup.Group("/beta/v1")
	//{
	//	NoAuthGroup := v1Group.Group("/")
	//	NoAuthGroup.Use(middleware.SessionMiddleware())
	//	JwtAuthGroup := v1Group.Group("/")
	//	JwtAuthGroup.Use(middleware.JwtMiddleware(jwt2, logger))
	//	KeyAuthGroup := v1Group.Group("/")
	//	KeyAuthGroup.Use(middleware.ApiKeyMiddleware(apiKeySvc, logger))
	//	KeyAuthGroupBeta := v1Beta.Group("/")
	//	KeyAuthGroupBeta.Use(middleware.ApiKeyMiddleware(apiKeySvc, logger))
	//	{
	//		// 无需鉴权
	//		NoAuthGroup.GET("/oauth2/linuxDo", oauth2Handler.LinuxDoLogin)
	//		NoAuthGroup.GET("/oauth2/linuxDo/callback", oauth2Handler.LinuxDoCallback)
	//		NoAuthGroup.GET("/oauth2/session", oauth2Handler.GetUserInfo)
	//		NoAuthGroup.POST("/auth/token/refresh", authHandler.GetNewAccessToken)
	//
	//		NoAuthGroup.POST("/system/email", sysConfigHandler.SetEmailConfig)
	//		NoAuthGroup.GET("/system/email", sysConfigHandler.GetEmailConfig)
	//		NoAuthGroup.GET("/system/email/health", sysConfigHandler.IsEmailServiceAvailable)
	//	}
	//	{
	//		// JWT需要鉴权
	//		JwtAuthGroup.GET("/key", apiKeyHandler.GetApiKey)
	//		JwtAuthGroup.PUT("/key", apiKeyHandler.ResetApiKey)
	//		JwtAuthGroup.GET("/user", userHandler.GetUser)
	//		JwtAuthGroup.GET("/reqLog/list", requestLogHandler.GetRequestLogs)
	//		JwtAuthGroup.GET("/reqLog/ranking", requestLogHandler.GetRequestLogRanking)
	//	}
	//	{
	//		// Key需要鉴权
	//		KeyAuthGroup.POST("/chat/completions", oaiHandler.ChatCompletions)
	//		KeyAuthGroup.POST("/completions", oaiHandler.Completions)
	//		KeyAuthGroup.GET("/models", oaiHandler.Models)
	//		KeyAuthGroup.POST("/embeddings", oaiHandler.Embeddings)
	//		KeyAuthGroup.POST("/audio/speech", oaiHandler.AudioSpeech)
	//		KeyAuthGroup.POST("/audio/transcriptions", oaiHandler.AudioTranscription)
	//		KeyAuthGroup.POST("/audio/translations", oaiHandler.AudioTranslation)
	//		KeyAuthGroup.POST("/images/generations", oaiHandler.ImageGeneration)
	//		KeyAuthGroup.POST("/images/edits", oaiHandler.ImageEdit)
	//		KeyAuthGroup.POST("/images/variations", oaiHandler.ImageVariation)
	//		KeyAuthGroupBeta.POST("/chat/completions", oaiHandler.ChatCompletionsByBytes)
	//		KeyAuthGroupBeta.POST("/completions", oaiHandler.CompletionsByBytes)
	//		KeyAuthGroupBeta.GET("/models", oaiHandler.Models)
	//		KeyAuthGroupBeta.POST("/embeddings", oaiHandler.EmbeddingsByBytes)
	//		KeyAuthGroupBeta.POST("/audio/speech", oaiHandler.AudioSpeechByBytes)
	//		KeyAuthGroupBeta.POST("/audio/transcriptions", oaiHandler.AudioTranscription)
	//		KeyAuthGroupBeta.POST("/audio/translations", oaiHandler.AudioTranslation)
	//		KeyAuthGroupBeta.POST("/images/generations", oaiHandler.ImageGenerationByBytes)
	//		KeyAuthGroupBeta.POST("/images/edits", oaiHandler.ImageEdit)
	//		KeyAuthGroupBeta.POST("/images/variations", oaiHandler.ImageVariation)
	//	}
	//}

	return s
}
