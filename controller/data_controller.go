package controller

import (
	"adorable-star/crawler"
	"adorable-star/model"
	"adorable-star/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DataController struct {
	s *service.DataService
}

func NewDataController(s *service.DataService) *DataController {
	return &DataController{s}
}

func (c *DataController) JupiterData(ctx *gin.Context) {
	// TODO: Change test codes into production code
	var params model.JupiterData
	ctx.ShouldBind(&params)

	// Crawl data asynchronously
	go crawler.FetchData(params.Account, params.Password)

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Crawler work in progress",
		"data": nil,
	})
}
