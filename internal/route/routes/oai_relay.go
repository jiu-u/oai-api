package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jiu-u/oai-api/internal/handler"
	"github.com/jiu-u/oai-api/internal/middleware"
	"github.com/jiu-u/oai-api/internal/service"
	"github.com/jiu-u/oai-api/pkg/log"
)

func SetupOaiRoutes(
	v1 *gin.RouterGroup,
	v1beta *gin.RouterGroup,
	apiKeyHandler *handler.ApiKeyHandler,
	oaiHandler *handler.OAIHandler,
	apiKeySvc service.ApiKeyService,
	logger *log.Logger,
) {
	r := v1.Group("/")
	r2 := v1beta.Group("/")
	keyAuthMiddleware := middleware.ApiKeyMiddleware(apiKeySvc, logger)
	r.Use(keyAuthMiddleware)
	r2.Use(keyAuthMiddleware)
	// 注册中间件
	{
		r.POST("/chat/completions", oaiHandler.ChatCompletions)
		r.POST("/completions", oaiHandler.Completions)
		r.GET("/models", oaiHandler.Models)
		r.POST("/embeddings", oaiHandler.Embeddings)
		r.POST("/audio/speech", oaiHandler.AudioSpeech)
		r.POST("/audio/transcriptions", oaiHandler.AudioTranscription)
		r.POST("/audio/translations", oaiHandler.AudioTranslation)
		r.POST("/images/generations", oaiHandler.ImageGeneration)
		r.POST("/images/edits", oaiHandler.ImageEdit)
		r.POST("/images/variations", oaiHandler.ImageVariation)
		r2.POST("/chat/completions", oaiHandler.ChatCompletionsByBytes)
		r2.POST("/completions", oaiHandler.CompletionsByBytes)
		r2.GET("/models", oaiHandler.Models)
		r2.POST("/embeddings", oaiHandler.EmbeddingsByBytes)
		r2.POST("/audio/speech", oaiHandler.AudioSpeechByBytes)
		r2.POST("/audio/transcriptions", oaiHandler.AudioTranscription)
		r2.POST("/audio/translations", oaiHandler.AudioTranslation)
		r2.POST("/images/generations", oaiHandler.ImageGenerationByBytes)
		r2.POST("/images/edits", oaiHandler.ImageEdit)
		r2.POST("/images/variations", oaiHandler.ImageVariation)
	}
}
