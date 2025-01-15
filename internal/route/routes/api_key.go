package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jiu-u/oai-api/internal/handler"
	"github.com/jiu-u/oai-api/internal/middleware"
	"github.com/jiu-u/oai-api/pkg/jwt"
	"github.com/jiu-u/oai-api/pkg/log"
)

func SetupApiKeyRoutes(
	v1 *gin.RouterGroup,
	apiKeyHandler *handler.ApiKeyHandler,
	jwtJWT *jwt.JWT,
	logger *log.Logger,
) {
	// 中间件
	keyGroup := v1.Group("/key")
	keyGroup.Use(middleware.JwtMiddleware(jwtJWT, logger))
	{
		keyGroup.GET("", apiKeyHandler.GetApiKey)
		keyGroup.POST("/create", apiKeyHandler.ResetApiKey)
		keyGroup.POST("/reset", apiKeyHandler.ResetApiKey)
	}

}
