package routes

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/jiu-u/oai-api/api/v1"
)

func ImplementHandle(c *gin.Context) {
	v1.HandleSuccess(c, gin.H{"message": "This feature is not implemented yet."})
}
