package controller

import (
	"adorable-star/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var GPT = &GptController{service.GPT}

type GptController struct {
	s *service.GptService
}

func (c *GptController) Conversation(ctx *gin.Context) {
	convId, response := c.s.Conversation("Who are you?")

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  response,
		"data": convId,
	})
}
