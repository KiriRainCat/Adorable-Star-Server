package controller

import (
	"adorable-star/internal/dao"
	"adorable-star/internal/pkg/config"
	"adorable-star/internal/pkg/crawler"
	"adorable-star/internal/pkg/util"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var Admin = &AdminController{}

type AdminController struct {
}

func (c *AdminController) SwitchBrowser(ctx *gin.Context) {
	// Parse params
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}

	crawler.SwitchBrowser(i)

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": nil,
	})
}

func (c *AdminController) GetCrawlerLog(ctx *gin.Context) {
	bytes, err := os.ReadFile(util.GetCwd() + "/storage/log/crawler.log")
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
		"data": string(bytes),
	})
}

func (c *AdminController) GetDbLog(ctx *gin.Context) {
	bytes, err := os.ReadFile(util.GetCwd() + "/storage/log/db.log")
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
		"data": string(bytes),
	})
}

func (c *AdminController) GetCrawlerLoad(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": strconv.Itoa(crawler.PagePoolLoad) + " / " +
			strconv.Itoa(config.Config.Crawler.MaxParallel) + "|" +
			strconv.Itoa(crawler.PendingTaskCount),
	})
}

func (c *AdminController) GetUsers(ctx *gin.Context) {
	users, err := dao.User.GetUsers()
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
		"data": users,
	})
}

func (c *AdminController) UpdateUserStatus(ctx *gin.Context) {
	// Parse params
	rawUid, rawStatus := ctx.Param("id"), ctx.Param("status")
	if rawUid == "" || rawStatus == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}
	uid, err1 := strconv.Atoi(rawUid)
	status, err2 := strconv.Atoi(rawStatus)
	if err1 != nil || err2 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}

	err := dao.User.UpdateStatus(uid, status)
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

func (c *AdminController) ChangeUserPassword(ctx *gin.Context) {
	type json struct {
		UID      int    `json:"uid" binding:"required"`
		Password string `json:"password" binding:"required"`
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

	// Update user password
	encodedPwd, err := bcrypt.GenerateFromPassword([]byte(data.Password+config.Config.Server.EncryptSalt), bcrypt.MinCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "服务器内部发生错误，请联系开发者",
			"data": nil,
		})
		return
	}

	err = dao.User.UpdatePassword(data.UID, string(encodedPwd))
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
