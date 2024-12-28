package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jiu-u/oai-api/pkg/jwt"
	"github.com/jiu-u/oai-api/pkg/log"
	"go.uber.org/zap"
	"net/http"
)

func JwtMiddleware(jwt *jwt.JWT, logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//tokenString := ctx.Request.Header.Get("Authorization")
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			logger.WithContext(ctx).Warn("No token", zap.Any("data", map[string]interface{}{
				"url":    ctx.Request.URL,
				"params": ctx.Params,
			}))
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token is empty"})
			return
		}

		claims, err := jwt.ParseAccessToken(tokenString, "Bearer ")
		if err != nil {
			logger.WithContext(ctx).Error("token error", zap.Any("data", map[string]interface{}{
				"url":    ctx.Request.URL,
				"params": ctx.Params,
			}), zap.Error(err))
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.Set("claims", claims)
		ctx.Next()
	}
}
