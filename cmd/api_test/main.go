package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jiu-u/oai-api/internal/middleware"
	"github.com/jiu-u/oai-api/pkg/adapter/provider"
)

func main() {
	p, err := NewProvider(provider.Config{
		Type:     "openai",
		EndPoint: "https://xxx.ai.com",
		APIKey:   "sk-1233213131321313",
	})
	if err != nil {
		panic(err)
	}
	oaiHandler := NewOpenAIHandler(p)
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	v1 := r.Group("/v1")
	{
		v1.POST("/chat/completions", oaiHandler.ChatCompletions)
	}
	bytesGroup := r.Group("/bytes/v1")
	{
		bytesGroup.POST("/chat/completions", oaiHandler.ChatCompletionsByBytes)
	}
	r.Run(":8080")
}
