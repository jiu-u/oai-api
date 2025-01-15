package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jiu-u/oai-api/internal/handler"
)

func SetupVerificationRoutes(v1 *gin.RouterGroup, h *handler.VerificationHandler) {
	vg := v1.Group("/verification")
	{
		vg.POST("/code/sms/send", ImplementHandle)
		vg.POST("/code/email/send", h.SetVerificationCode2Email)
	}
}
