package controller

import (
	"adorable-star/internal/pkg/crawler"
	"adorable-star/internal/service"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var GPT = &GptController{service.GPT, time.Now().Unix()}

type GptController struct {
	s            *service.GptService
	requested_at int64
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

	// Check time elapsed from last request to decide whether required to reload the page
	now := time.Now().Unix()
	if now-c.requested_at > 300 {
		c.requested_at = now
		crawler.GetGptAccessToken()
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
