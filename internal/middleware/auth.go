package middleware

import (
	"adorable-star/internal/pkg/config"
	"adorable-star/internal/service"
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
	if ctx.Request.Header.Get("Authorization") != config.Config.Server.RequestAuth {
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
		strings.Contains(ctx.Request.URL.String(), "register") ||
		strings.Contains(ctx.Request.URL.String(), "complete-info") {
		ctx.Next()
		return
	}

	// Authenticate Request Header
	claims, err := m.s.VerifyToken(ctx)
	if err != nil {
		println(err.Error())
		if err.Error() == "internalErr" {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "服务器内部发生错误，请联系开发者",
				"data": nil,
			})
			return
		}
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "用户 Token 验证不通过",
			"data": nil,
		})
		ctx.Abort()
		return
	}

	ctx.Set("uid", claims.UID)
	ctx.Next()
}

func (m *AuthMiddleware) AuthenticateAdmin(ctx *gin.Context) {
	// Let PING api to pass
	if strings.Contains(ctx.Request.URL.String(), "ping") {
		ctx.Next()
		return
	}

	// Authenticate Request Header
	if ctx.Request.Header.Get("Admin") != config.Config.Server.AdminAuth {
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
