package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SessionMiddleware() gin.HandlerFunc {
	store := cookie.NewStore([]byte("secret"))
	return sessions.Sessions("session", store)
}
