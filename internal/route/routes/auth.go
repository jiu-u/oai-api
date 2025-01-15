package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jiu-u/oai-api/internal/handler"
	"github.com/jiu-u/oai-api/internal/middleware"
	"github.com/jiu-u/oai-api/pkg/jwt"
	"github.com/jiu-u/oai-api/pkg/log"
)

func SetupAuthRoutes(
	v1 *gin.RouterGroup,
	authHandler *handler.AuthHandler,
	sysConfigHandler *handler.SystemConfigHandler,
	userHandler *handler.UserHandler,
	jwtJWT *jwt.JWT,
	logger *log.Logger,
) {
	// 中间件
	authGroup := v1.Group("/auth")
	NoAuthGroup := authGroup.Group("/")
	needAuthGroup := authGroup.Group("/")
	needAuthGroup.Use(middleware.JwtMiddleware(jwtJWT, logger))
	{
		needAuthGroup.POST("/logout", ImplementHandle)
		needAuthGroup.GET("/current-user", userHandler.GetCurrentUser)

		NoAuthGroup.POST("/login", authHandler.Login)
		NoAuthGroup.POST("/register", authHandler.Register)
		NoAuthGroup.POST("/access-token", authHandler.GetNewAccessToken)
		NoAuthGroup.GET("/session", authHandler.LoginBySessionId)
		// 检查是否可用
		NoAuthGroup.GET("/login/oauth2/linux-do", sysConfigHandler.IsLinuxDoOAuthServiceAvailable)
		NoAuthGroup.GET("/login/oauth2/linux-do/redirect", authHandler.LinuxDoLogin)
		NoAuthGroup.GET("/oauth2/linux-do/callback", authHandler.LinuxDoCallBack)
		// 检查是否可用
		NoAuthGroup.GET("/login/oauth2/github", sysConfigHandler.IsGithubOAuthServiceAvailable)
		NoAuthGroup.GET("/login/oauth2/github/redirect", authHandler.GithubLogin)
		NoAuthGroup.GET("/oauth2/github/callback", authHandler.GithubCallBack)
	}
	//r.POST("/")

}
