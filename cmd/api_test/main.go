package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jiu-u/oai-api/internal/middleware"
	"github.com/jiu-u/oai-api/pkg/adapter/provider"
)

func main() {
	p, err := NewProvider(provider.Config{
		Type:     "openai",
		EndPoint: "https://api.ai.com",
		APIKey:   "sk-12332131321131321",
	})
	if err != nil {
		panic(err)
	}
	oaiHandler := NewOpenAIHandler(p)
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	v1 := r.Group("/v1")
	fmt.Println("------------------------")
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
	bytesGroup := r.Group("/bytes/v1")
	{
		bytesGroup.POST("/chat/completions", oaiHandler.ChatCompletionsByBytes)
		bytesGroup.POST("/completions", oaiHandler.CompletionsByBytes)
		bytesGroup.GET("/models", oaiHandler.Models)
		bytesGroup.POST("/embeddings", oaiHandler.EmbeddingsByBytes)
		bytesGroup.POST("/audio/speech", oaiHandler.AudioSpeechByBytes)
		bytesGroup.POST("/audio/transcriptions", oaiHandler.AudioTranscription)
		bytesGroup.POST("/audio/translations", oaiHandler.AudioTranslation)
		bytesGroup.POST("/images/generations", oaiHandler.ImageGenerationByBytes)
		bytesGroup.POST("/images/edits", oaiHandler.ImageEdit)
		bytesGroup.POST("/images/variations", oaiHandler.ImageVariation)
	}
	r.Run(":8080")
}
