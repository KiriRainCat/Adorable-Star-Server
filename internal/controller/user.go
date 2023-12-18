package controller

import (
	"adorable-star/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var User = &AuthController{service.User}

type AuthController struct {
	s *service.UserService
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
	token, user, isReturningUser, err := c.s.Login(data.Name, data.Password)
	if err != nil {
		if err.Error() == "internalErr" {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "服务器内部发生错误，请联系开发者",
				"data": nil,
			})
			return
		}
		if err.Error() == "userJupiterDataNotFound" {
			ctx.JSON(http.StatusPreconditionRequired, gin.H{
				"code": http.StatusPreconditionRequired,
				"msg":  "需要用户添加 Jupiter 数据",
				"data": user.ID,
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
		"msg":  "success",
		"data": token + "|" + strconv.Itoa(user.Status) + "|" + strconv.FormatBool(isReturningUser),
	})
}

func (c *AuthController) ValidationCode(ctx *gin.Context) {
	userMail := ctx.Param("email")

	// When queries are empty
	if userMail == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "账号或密码错误",
			"data": nil,
		})
		return
	}

	// Send validation code
	err := c.s.SendValidationCode(userMail)
	if err != nil {
		if strings.Contains(err.Error(), "internalErr") {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "服务器内部发生错误，请联系开发者",
				"data": nil,
			})
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
		"msg":  "success",
		"data": nil,
	})
}

func (c *AuthController) Register(ctx *gin.Context) {
	type json struct {
		Email          string `json:"email,omitempty" binding:"required"`
		ValidationCode string `json:"validation_code,omitempty"`
		Username       string `json:"username,omitempty" binding:"required"`
		Password       string `json:"password,omitempty" binding:"required"`
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
	if err := c.s.Register(data.Email, data.ValidationCode, data.Username, data.Password); err != nil {
		if err.Error() == "internalErr" {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "服务器内部发生错误，请联系开发者",
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
		"msg":  "success",
		"data": nil,
	})
}

func (c *AuthController) ChangePassword(ctx *gin.Context) {
	type json struct {
		Password    string `json:"password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	// When queries are empty
	var data json
	if ctx.ShouldBind(&data) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}

	// Change password
	err := c.s.ChangePassword(ctx.GetInt("uid"), data.Password, data.NewPassword)
	if err != nil {
		if err.Error() == "internalErr" {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "服务器内部发生错误，请联系开发者",
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
		"msg":  "success",
		"data": nil,
	})
}

func (c *AuthController) ChangeCfbp(ctx *gin.Context) {
	cfbp := ctx.Param("cfbp")

	// When queries are empty
	if cfbp == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}

	// Update cfbp
	err := c.s.ChangeCfbp(ctx.GetInt("uid"), cfbp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "服务器内部发生错误，请联系开发者",
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": nil,
	})
}

func (c *AuthController) CompleteInfo(ctx *gin.Context) {
	type json struct {
		UID      int    `json:"uid" binding:"required"`
		Account  string `json:"account" binding:"required"`
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

	// Check and insert user's Jupiter data
	if err := c.s.CompleteInfo(data.UID, data.Account, data.Password); err != nil {
		if err.Error() == "internalErr" {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "服务器内部发生错误，请联系开发者",
				"data": nil,
			})
			return
		}
		ctx.JSON(http.StatusExpectationFailed, gin.H{
			"code": http.StatusExpectationFailed,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": nil,
	})
}
