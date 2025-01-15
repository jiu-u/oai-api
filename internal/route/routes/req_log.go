package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jiu-u/oai-api/internal/handler"
	"github.com/jiu-u/oai-api/internal/middleware"
	"github.com/jiu-u/oai-api/pkg/jwt"
	"github.com/jiu-u/oai-api/pkg/log"
)

func SetupOaiReqLogRoutes(
	v1 *gin.RouterGroup,
	requestLogHandler *handler.RequestLogHandler,
	jwtJWT *jwt.JWT,
	logger *log.Logger,
) {
	logsGroup := v1.Group("/oai-logs")
	logsGroup.Use(middleware.JwtMiddleware(jwtJWT, logger))
	{
		logsGroup.GET("", requestLogHandler.GetRequestLogs)
		// 查询某个用户调用请求日志
		logsGroup.GET("/users/:userId", requestLogHandler.GetUserRequestLogs)
		// 用户调用次数排行
		logsGroup.GET("/users-ranking", requestLogHandler.GetRequestLogsUserRanking)
		// 模型调用次数排行
		logsGroup.GET("/models-ranking", requestLogHandler.GetRequestLogsModelRanking)
	}

}
