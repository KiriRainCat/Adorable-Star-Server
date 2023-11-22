package router

import (
	"adorable-star/internal/controller"
	"adorable-star/internal/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup) {
	g := r.Group("/user")

	// Deps
	c := controller.User
	m := middleware.Auth

	// Routes
	g.POST("/login", c.Login)
	g.POST("/complete-info", c.CompleteInfo)
	g.POST("/register", m.AuthenticateAdmin, c.Register)
	g.PUT("/password", c.ChangePassword)
	g.PUT("/cfbp/:cfbp", c.ChangeCfbp)
}
