package controller

import (
	"adorable-star/internal/service"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

var GPT = &GptController{service.GPT}

type GptController struct {
	s *service.GptService
}

func (c *GptController) Conversation(ctx *gin.Context) {
	type json struct {
		ConversationId string              `json:"conversation_id,omitempty"`
		Messages       []map[string]string `json:"messages" binding:"required"`
	}

	// Get query and check if it's empty
	var data json
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数错误",
		})
		return
	}

	stream := make(chan string, 10)
	go func() {
		defer close(stream)
		c.s.Conversation(stream, data.ConversationId, data.Messages)
	}()

	ctx.Stream(func(w io.Writer) bool {
		if text, ok := <-stream; ok {
			ctx.SSEvent("message", text)
			return true
		} else {
			return false
		}
	})
}
