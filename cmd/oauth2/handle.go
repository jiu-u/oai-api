package oauth2

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type oauthData struct {
	Code string `json:"code" form:"code"`
	//proto.Stat
}

func LinuxDoAuthHandle(ctx *gin.Context) {
	fmt.Println("code", ctx.Query("code"))
	fmt.Println("state", ctx.Query("state"))
	//ctx.Request.SetBasicAuth().
	
}
