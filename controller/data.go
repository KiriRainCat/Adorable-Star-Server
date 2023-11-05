package controller

import (
	"adorable-star/service/crawler"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var Data = &DataController{}

type DataController struct{}

func (c *DataController) JupiterData(ctx *gin.Context) {
	uid, _ := strconv.Atoi(ctx.Query("uid"))

	// Crawl data asynchronously
	go crawler.FetchData(uid)

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Crawler work in progress",
		"data": nil,
	})
}
