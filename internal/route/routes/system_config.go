package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jiu-u/oai-api/internal/handler"
	"github.com/jiu-u/oai-api/internal/middleware"
	"github.com/jiu-u/oai-api/pkg/jwt"
	"github.com/jiu-u/oai-api/pkg/log"
)

func SetupSystemConfigRoutes(
	v1 *gin.RouterGroup,
	sysConfigHandler *handler.SystemConfigHandler,
	jwtJWT *jwt.JWT,
	logger *log.Logger,
) {
	systemGroup := v1.Group("/system")
	needAuthGroup := systemGroup.Group("/")
	needAuthGroup.Use(middleware.JwtMiddleware(jwtJWT, logger))
	// no auth
	{
		systemGroup.GET("/email/health", sysConfigHandler.IsEmailServiceAvailable)
		systemGroup.GET("/linux-do/health", sysConfigHandler.IsLinuxDoOAuthServiceAvailable)
		systemGroup.GET("/github/health", sysConfigHandler.IsGithubOAuthServiceAvailable)
		systemGroup.GET("/register", sysConfigHandler.GetRegisterConfig)
		systemGroup.GET("/model", sysConfigHandler.GetModelConfig)
	}
	{
		// need auth
		needAuthGroup.POST("/register", sysConfigHandler.SetRegisterConfig)
		needAuthGroup.POST("/email", sysConfigHandler.SetEmailConfig)
		needAuthGroup.GET("/email", sysConfigHandler.GetEmailConfig)
		needAuthGroup.POST("/linux-do", sysConfigHandler.SetLinuxDoOAuthConfig)
		needAuthGroup.GET("/linux-do", sysConfigHandler.GetLinuxDoOAuthConfig)
		needAuthGroup.POST("/github", sysConfigHandler.SetGithubOAuthConfig)
		needAuthGroup.GET("/github", sysConfigHandler.GetGithubOAuthConfig)
		needAuthGroup.POST("/model", sysConfigHandler.SetModelConfig)
	}
}
