package server

import (
	"embed"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jiu-u/oai-api/internal/handler"
	"github.com/jiu-u/oai-api/internal/middleware"
	"github.com/jiu-u/oai-api/internal/service"
	"github.com/jiu-u/oai-api/pkg/config"
	"github.com/jiu-u/oai-api/pkg/jwt"
	"github.com/jiu-u/oai-api/pkg/log"
	"github.com/jiu-u/oai-api/pkg/server/http"
	"io/fs"
	stdhttp "net/http"
)

type embedFileSystem struct {
	stdhttp.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	if err != nil {
		return false
	}
	return true
}

func EmbedFolder(fsEmbed embed.FS, targetPath string) static.ServeFileSystem {
	efs, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: stdhttp.FS(efs),
	}
}

////go:embed web/dist
//var buildFS embed.FS
//
////go:embed web/dist/index.html
//var indexPage []byte

func NewHTTPServer(
	logger *log.Logger,
	cfg *config.Config,
	jwt2 *jwt.JWT,
	oaiHandler *handler.OAIHandler,
	oauth2Handler *handler.OAuth2Handler,
	authHandler *handler.AuthHandler,
	apiKeySvc service.ApiKeyService,
	apiKeyHandler *handler.ApiKeyHandler,
	userHandler *handler.UserHandler,
	requestLogHandler *handler.RequestLogHandler,
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
	)
	s.Static("/assets", "./web/dist/assets")

	//s.Use(static.Serve("/", EmbedFolder(buildFS, "web/dist")))
	//s.GET("/", func(ctx *gin.Context) {
	//	apiV1.HandleSuccess(ctx, map[string]interface{}{
	//		":)": "Hello",
	//	})
	//})
	s.GET("/", func(ctx *gin.Context) {
		ctx.File("./web/dist/index.html")
	})
	apiGroup := s.Group("/api")
	v1Group := apiGroup.Group("/v1")
	v1Beta := apiGroup.Group("/beta/v1")
	{
		NoAuthGroup := v1Group.Group("/")
		NoAuthGroup.Use(middleware.SessionMiddleware())
		JwtAuthGroup := v1Group.Group("/")
		JwtAuthGroup.Use(middleware.JwtMiddleware(jwt2, logger))
		KeyAuthGroup := v1Group.Group("/")
		KeyAuthGroup.Use(middleware.ApiKeyMiddleware(apiKeySvc, logger))
		KeyAuthGroupBeta := v1Beta.Group("/")
		KeyAuthGroupBeta.Use(middleware.ApiKeyMiddleware(apiKeySvc, logger))
		{
			// 无需鉴权
			NoAuthGroup.GET("/oauth2/linuxDo", oauth2Handler.LinuxDoLogin)
			NoAuthGroup.GET("/oauth2/linuxDo/callback", oauth2Handler.LinuxDoCallback)
			NoAuthGroup.GET("/oauth2/session", oauth2Handler.GetUserInfo)
			NoAuthGroup.POST("/auth/token/refresh", authHandler.GetNewAccessToken)
		}
		{
			// JWT需要鉴权
			JwtAuthGroup.GET("/key", apiKeyHandler.GetApiKey)
			JwtAuthGroup.PUT("/key", apiKeyHandler.ResetApiKey)
			JwtAuthGroup.GET("/user", userHandler.GetUser)
			JwtAuthGroup.GET("/reqLog/list", requestLogHandler.GetRequestLogs)
			JwtAuthGroup.GET("/reqLog/ranking", requestLogHandler.GetRequestLogRanking)
		}
		{
			// Key需要鉴权
			KeyAuthGroup.POST("/chat/completions", oaiHandler.ChatCompletions)
			KeyAuthGroup.POST("/completions", oaiHandler.Completions)
			KeyAuthGroup.GET("/models", oaiHandler.Models)
			KeyAuthGroup.POST("/embeddings", oaiHandler.Embeddings)
			KeyAuthGroup.POST("/audio/speech", oaiHandler.AudioSpeech)
			KeyAuthGroup.POST("/audio/transcriptions", oaiHandler.AudioTranscription)
			KeyAuthGroup.POST("/audio/translations", oaiHandler.AudioTranslation)
			KeyAuthGroup.POST("/images/generations", oaiHandler.ImageGeneration)
			KeyAuthGroup.POST("/images/edits", oaiHandler.ImageEdit)
			KeyAuthGroup.POST("/images/variations", oaiHandler.ImageVariation)
			KeyAuthGroupBeta.POST("/chat/completions", oaiHandler.ChatCompletionsByBytes)
			KeyAuthGroupBeta.POST("/completions", oaiHandler.CompletionsByBytes)
			KeyAuthGroupBeta.GET("/models", oaiHandler.Models)
			KeyAuthGroupBeta.POST("/embeddings", oaiHandler.EmbeddingsByBytes)
			KeyAuthGroupBeta.POST("/audio/speech", oaiHandler.AudioSpeechByBytes)
			KeyAuthGroupBeta.POST("/audio/transcriptions", oaiHandler.AudioTranscription)
			KeyAuthGroupBeta.POST("/audio/translations", oaiHandler.AudioTranslation)
			KeyAuthGroupBeta.POST("/images/generations", oaiHandler.ImageGenerationByBytes)
			KeyAuthGroupBeta.POST("/images/edits", oaiHandler.ImageEdit)
			KeyAuthGroupBeta.POST("/images/variations", oaiHandler.ImageVariation)
		}
	}
	//
	//v1 := s.Group("/v1")
	//{
	//	v1.POST("/chat/completions", oaiHandler.ChatCompletions)
	//	v1.POST("/completions", oaiHandler.Completions)
	//	v1.GET("/models", oaiHandler.Models)
	//	v1.POST("/embeddings", oaiHandler.Embeddings)
	//	v1.POST("/audio/speech", oaiHandler.AudioSpeech)
	//	v1.POST("/audio/transcriptions", oaiHandler.AudioTranscription)
	//	v1.POST("/audio/translations", oaiHandler.AudioTranslation)
	//	v1.POST("/images/generations", oaiHandler.ImageGeneration)
	//	v1.POST("/images/edits", oaiHandler.ImageEdit)
	//	v1.POST("/images/variations", oaiHandler.ImageVariation)
	//
	//}
	//bytesV1 := s.Group("/bytes/v1")
	//{
	//	bytesV1.POST("/chat/completions", oaiHandler.ChatCompletionsByBytes)
	//	bytesV1.POST("/completions", oaiHandler.CompletionsByBytes)
	//	bytesV1.GET("/models", oaiHandler.Models)
	//	bytesV1.POST("/embeddings", oaiHandler.EmbeddingsByBytes)
	//	bytesV1.POST("/audio/speech", oaiHandler.AudioSpeechByBytes)
	//	bytesV1.POST("/audio/transcriptions", oaiHandler.AudioTranscription)
	//	bytesV1.POST("/audio/translations", oaiHandler.AudioTranslation)
	//	bytesV1.POST("/images/generations", oaiHandler.ImageGenerationByBytes)
	//	bytesV1.POST("/images/edits", oaiHandler.ImageEdit)
	//	bytesV1.POST("/images/variations", oaiHandler.ImageVariation)
	//}

	return s
}
