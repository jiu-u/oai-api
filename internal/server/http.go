package server

import (
	"github.com/gin-gonic/gin"
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/handler"
	"github.com/jiu-u/oai-api/internal/middleware"
	"github.com/jiu-u/oai-api/pkg/config"
	"github.com/jiu-u/oai-api/pkg/log"
	"github.com/jiu-u/oai-api/pkg/server/http"
)

func NewHTTPServer(
	logger *log.Logger,
	cfg *config.Config,
	oaiHandler *handler.OAIHandler,
) *http.Server {
	gin.SetMode(gin.DebugMode)
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
	s.GET("/", func(ctx *gin.Context) {
		apiV1.HandleSuccess(ctx, map[string]interface{}{
			":)": "Hello",
		})
	})

	v1 := s.Group("/v1")
	{
		v1.POST("/chat/completions", oaiHandler.ChatCompletions)
		v1.POST("/completions", oaiHandler.Completions)
		v1.GET("/models", oaiHandler.Models)
		v1.POST("/embeddings", oaiHandler.Embeddings)
		v1.POST("/audio/speech", oaiHandler.AudioSpeech)
		v1.POST("/audio/transcriptions", oaiHandler.AudioTranscription)
		v1.POST("/audio/translations", oaiHandler.AudioTranslation)
		v1.POST("/images/generations", oaiHandler.ImageGeneration)
		v1.POST("/images/edits", oaiHandler.ImageEdit)
		v1.POST("/images/variations", oaiHandler.ImageVariation)

	}
	bytesV1 := s.Group("/bytes/v1")
	{
		bytesV1.POST("/chat/completions", oaiHandler.ChatCompletionsByBytes)
		bytesV1.POST("/completions", oaiHandler.CompletionsByBytes)
		bytesV1.GET("/models", oaiHandler.Models)
		bytesV1.POST("/embeddings", oaiHandler.EmbeddingsByBytes)
		bytesV1.POST("/audio/speech", oaiHandler.AudioSpeechByBytes)
		bytesV1.POST("/audio/transcriptions", oaiHandler.AudioTranscription)
		bytesV1.POST("/audio/translations", oaiHandler.AudioTranslation)
		bytesV1.POST("/images/generations", oaiHandler.ImageGenerationByBytes)
		bytesV1.POST("/images/edits", oaiHandler.ImageEdit)
		bytesV1.POST("/images/variations", oaiHandler.ImageVariation)
	}

	return s
}
