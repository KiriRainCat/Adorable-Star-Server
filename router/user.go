package router

import (
	"adorable-star/internal/controller"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup) {
	g := r.Group("/user")

	// Deps
	c := controller.User

	// Routes
	g.POST("/login", c.Login)
	g.POST("/complete-info", c.CompleteInfo)
	g.POST("/validation-code/:email", c.ValidationCode)
	g.POST("/register", c.Register)
	g.PUT("/password", c.ChangePassword)
	g.PUT("/cfbp/:cfbp", c.ChangeCfbp)
}
