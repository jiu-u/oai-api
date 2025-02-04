package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jiu-u/oai-api/constant"
	"github.com/jiu-u/oai-api/pkg/log"
	"github.com/lithammer/shortuuid/v4"
	"go.uber.org/zap"
)

func TraceMiddleware(logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		Ip := ctx.ClientIP()
		ctx.Set(constant.ClientIPKey, Ip)
		traceId := shortuuid.New()
		logger.WithValue(ctx, zap.String("traceId", traceId), zap.String("type", "request"))
		ctx.Next()
	}
}
