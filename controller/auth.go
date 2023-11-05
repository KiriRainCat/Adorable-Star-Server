package controller

import (
	"adorable-star/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var Auth = &AuthController{service.Auth}

type AuthController struct {
	s *service.AuthService
}

func (c *AuthController) Login(ctx *gin.Context) {
	type json struct {
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// When queries are empty
	var data json
	if ctx.ShouldBind(&data) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "账号或密码错误",
			"data": nil,
		})
		return
	}

	// Login
	token, err := c.s.Login(data.Name, data.Password)
	if err != nil {
		if err.Error() == "internalErr" {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "服务端内部发生错误",
				"data": nil,
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "登录成功",
		"data": token,
	})
}

func (c *AuthController) Logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "登出成功",
		"data": nil,
	})
}

func (c *AuthController) Register(ctx *gin.Context) {
	type json struct {
		Email    string `json:"email" binding:"required"`
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// When queries are empty
	var data json
	if ctx.ShouldBind(&data) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数不得为空",
			"data": nil,
		})
		return
	}

	// Register
	if err := c.s.Register(data.Email, data.Username, data.Password); err != nil {
		if err.Error() == "internalErr" {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "服务端内部发生错误",
				"data": nil,
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "注册成功",
		"data": nil,
	})
}
