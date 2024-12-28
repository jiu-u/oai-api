package handler

import (
	"github.com/gin-gonic/gin"
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/service"
	"github.com/jiu-u/oai-api/pkg/jwt"
)

type AuthHandler struct {
	*Handler
	jwt *jwt.JWT
	svc service.AuthService
}

func NewAuthHandler(handler *Handler, jwt *jwt.JWT, svc service.AuthService) *AuthHandler {
	return &AuthHandler{
		Handler: handler,
		jwt:     jwt,
		svc:     svc,
	}
}

func (h *AuthHandler) GetNewAccessToken(ctx *gin.Context) {
	req := new(apiV1.AccessTokenRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	claim, err := h.jwt.ParseRefreshToken(req.RefreshToken, "Bearer ")
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.svc.NewAccessToken(ctx, claim.UserId)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, resp)
}
