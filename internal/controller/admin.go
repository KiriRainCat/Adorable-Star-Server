package controller

import (
	"adorable-star/internal/pkg/config"
	"adorable-star/internal/pkg/crawler"
	"adorable-star/internal/pkg/util"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

var Admin = &AdminController{}

type AdminController struct {
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

func (c *AdminController) GetSqlLog(ctx *gin.Context) {
	bytes, err := os.ReadFile(util.GetCwd() + "/storage/log/sql.log")
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
			strconv.Itoa(crawler.TaskCount),
	})
}
