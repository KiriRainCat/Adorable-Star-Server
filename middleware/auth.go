package middleware

import (
	"adorable-star/config"
	"adorable-star/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var Auth = &AuthMiddleware{service.Token}

type AuthMiddleware struct {
	s *service.TokenService
}

func (m *AuthMiddleware) Authenticate(ctx *gin.Context) {
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

func (m *AuthMiddleware) AuthenticateUser(ctx *gin.Context) {
	// Let PING api to pass
	if strings.Contains(ctx.Request.URL.String(), "ping") ||
		strings.Contains(ctx.Request.URL.String(), "login") ||
		strings.Contains(ctx.Request.URL.String(), "register") {
		ctx.Next()
		return
	}

	// Authenticate Request Header
	if m.s.VerifyToken(ctx.Request.Header.Get("Token")) != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "用户 Token 验证不通过",
			"data": nil,
		})
		ctx.Abort()
		return
	}

	ctx.Next()
}

func (m *AuthMiddleware) AuthenticateAdmin(ctx *gin.Context) {
	// Let PING api to pass
	if strings.Contains(ctx.Request.URL.String(), "ping") {
		ctx.Next()
		return
	}

	// Authenticate Request Header
	if ctx.Request.Header.Get("Admin") != config.ADMIN_AUTH {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "管理员验证不通过",
			"data": nil,
		})
		ctx.Abort()
		return
	}

	ctx.Next()
}
