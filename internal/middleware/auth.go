package middleware

import (
	"adorable-star/internal/dao"
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
		strings.Contains(ctx.Request.URL.String(), "complete-info") ||
		strings.Contains(ctx.Request.URL.String(), "validation-code") ||
		strings.Contains(ctx.Request.URL.String(), "gpt") {
		ctx.Next()
		return
	}

	// Authenticate Request Header
	claims, err := m.s.VerifyToken(ctx)
	if err != nil {
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
	ctx.Set("status", claims.Status)

	// Set essential data
	data, err := dao.Jupiter.GetDataByUID(ctx.GetInt("uid"))
	if err != nil {
		ctx.Abort()
		return
	}
	ctx.Set("fetchedAt", data.FetchedAt)
	ctx.Set("gpa", data.GPA)

	ctx.Next()
}

// level 为等级，例如 10, 100 等，用于将用户 status 除以此值得出需要判断的 status
func (m *AuthMiddleware) AuthenticateUserLevel(level int, status int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Authenticate user status as level
		if ctx.GetInt("status")/level < status {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "权限不足",
				"data": nil,
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (m *AuthMiddleware) AuthenticateAdmin(ctx *gin.Context) {
	// Let PING api to pass
	if strings.Contains(ctx.Request.URL.String(), "ping") {
		ctx.Next()
		return
	}

	// Authenticate Request Header
	if ctx.Request.Header.Get("Admin") != config.Config.Server.AdminAuth || ctx.GetInt("status") != 900 {
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
