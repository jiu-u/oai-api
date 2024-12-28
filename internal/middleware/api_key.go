package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jiu-u/oai-api/internal/service"
	"github.com/jiu-u/oai-api/pkg/log"
	"net/http"
	"strings"
)

func ApiKeyMiddleware(apiKeySvc service.ApiKeyService, logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("Authorization")
		if apiKey == "" || !apiKeySvc.IsActiveApiKey(ctx, strings.TrimPrefix(apiKey, "Bearer ")) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "api key is invalid"})
			return
		}
		key := strings.TrimPrefix(apiKey, "Bearer ")
		ctx.Set("apiKey", key)
		ctx.Next()
	}
}
