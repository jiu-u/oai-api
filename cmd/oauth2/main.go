package oauth2

import (
	"github.com/gin-gonic/gin"
	"github.com/jiu-u/oai-api/internal/middleware"
	"net/http"
)

func main() {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"data": "good!",
		})
	})
	r.Run(":8888")
}
