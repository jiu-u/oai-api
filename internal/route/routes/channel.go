package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jiu-u/oai-api/internal/handler"
	"github.com/jiu-u/oai-api/internal/middleware"
	"github.com/jiu-u/oai-api/pkg/jwt"
	"github.com/jiu-u/oai-api/pkg/log"
)

func SetupChannelRoutes(
	v1 *gin.RouterGroup,
	channelHandler *handler.ChannelHandler,
	jwtJWT *jwt.JWT,
	logger *log.Logger,
) {
	// 中间件
	channelGroup := v1.Group("/channels")
	channelGroup.Use(middleware.JwtMiddleware(jwtJWT, logger))
	{
		channelGroup.GET("", channelHandler.GetChannels)
		channelGroup.GET("/:channelId", channelHandler.GetChannel)
		channelGroup.POST("", channelHandler.CreateChannel)
		channelGroup.PUT("/:channelId", channelHandler.UpdateChannel)
		channelGroup.PUT("/:channelId/status", channelHandler.UpdateChannelStatus)
		channelGroup.DELETE("/:channelId", channelHandler.DeleteChannel)
		channelGroup.POST("/:channelId/models/check", channelHandler.CheckModel)
		channelGroup.POST("/models/fetch", ImplementHandle)
		
		// 获取models
		//channelGroup.POST("/:channelId/models", ImplementHandle)
		// 设置models
		//channelGroup.PUT("/:channelId/models", ImplementHandle)
		//channelGroup.DELETE("/:channelId/models", ImplementHandle)
		//channelGroup.DELETE("/:channelId/models/:modelName", ImplementHandle)
	}

}
