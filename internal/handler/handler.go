package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jiu-u/oai-api/pkg/jwt"
	"github.com/jiu-u/oai-api/pkg/log"
)

type Handler struct {
	logger *log.Logger
}

func NewHandler(logger *log.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func GetUserIdFromCtx(ctx *gin.Context) uint64 {
	v, exists := ctx.Get("claims")
	if !exists {
		return 0
	}
	return v.(*jwt.MyCustomClaims).UserId
}

func GetUserRoleFromCtx(ctx *gin.Context) string {
	v, exists := ctx.Get("claims")
	if !exists {
		return ""
	}
	return v.(*jwt.MyCustomClaims).Role
}
