package middleware

import (
	"adorable-star/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AppAuthMiddleware struct{}

func (m *AppAuthMiddleware) Authenticate(ctx *gin.Context) {
	// Let PING api to pass
	if strings.Contains(ctx.Request.URL.String(), "ping") {
		ctx.Next()
		return
	}

	// Authenticate Request Header
	if ctx.Request.Header.Get("Authorization") != config.REQUEST_AUTH {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "客户端验证不通过",
			"data": nil,
		})
		ctx.Abort()
		return
	}

	ctx.Next()
}
