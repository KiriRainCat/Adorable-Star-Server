package controller

import (
	"adorable-star/config"
	"adorable-star/service"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func NewAuthController(s *service.UserService) *AuthController {
	return &AuthController{s}
}

type AuthController struct {
	s *service.UserService
}

func (c *AuthController) Login(ctx *gin.Context) {
	name := ctx.Query("name")
	password := ctx.Query("password") + config.ENCRYPT_SALT

	// When queries are empty
	if name == "" || password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "账号或密码错误",
			"data": nil,
		})
		return
	}

	// Find user from DB
	user, err := c.s.GetUserByUsernameOrEmail(name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "账户不存在",
			"data": nil,
		})
		return
	}

	// When pwd does not match
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "账号或密码错误",
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "登录成功",
		"data": fmt.Sprintf("%v>.<%v>.<%v", user.Status, name, user.CreatedAt.Unix()),
	})
}

func (c *AuthController) Logout(ctx *gin.Context) {
	token := strings.Split(ctx.Query("token"), ">.<")

	// When token length is not valid
	if len(token) < 3 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "非法访问",
			"data": nil,
		})
		return
	}

	// Find user from DB
	user, err := c.s.GetUserByUsernameOrEmail(token[1])
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "账户不存在",
			"data": nil,
		})
		return
	}

	// When token info mismatch with DB info
	if num, _ := strconv.Atoi(token[0]); num != user.Status || token[2] != strconv.FormatInt(user.CreatedAt.Unix(), 10) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "非法访问",
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "退出登录成功",
		"data": nil,
	})
}

func (c *AuthController) Register(ctx *gin.Context) {
	type json struct {
		Admin    string `json:"admin" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// When queries are empty
	var data json
	if ctx.ShouldBind(&data) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数不得为空 | 需要管理权限",
			"data": nil,
		})
		return
	}

	// Encrypt pwd
	encryptedPwd, err := bcrypt.GenerateFromPassword([]byte(data.Password+config.ENCRYPT_SALT), bcrypt.MinCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "服务端发生内部错误",
			"data": nil,
		})
		return
	}

	// Insert user to DB
	err = c.s.InsertUser(data.Email, data.Username, string(encryptedPwd[:]))
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "使用本邮箱或用户名的账户已经存在",
				"data": nil,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "服务端发生内部错误",
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
