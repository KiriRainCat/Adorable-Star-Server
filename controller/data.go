package controller

import (
	"adorable-star/service"
	"adorable-star/service/crawler"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var Data = &DataController{service.Data}

type DataController struct {
	s *service.DataService
}

func (c *DataController) JupiterData(ctx *gin.Context) {
	// TODO: Change test codes into production code
	uid, _ := strconv.Atoi(ctx.Query("uid"))

	// Crawl data asynchronously
	go crawler.FetchData(uid)

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Crawler work in progress",
		"data": nil,
	})
}
