package router

import (
	"adorable-star/internal/controller"

	"github.com/gin-gonic/gin"
)

func GptRoutes(r *gin.RouterGroup) {
	g := r.Group("/gpt")

	// Deps
	c := controller.GPT

	// Routes
	g.POST("/conversation", c.Conversation)
}
