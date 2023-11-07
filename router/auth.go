package router

import (
	"adorable-star/internal/controller"
	"adorable-star/internal/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup) {
	g := r.Group("/auth")

	// Deps
	c := controller.Auth
	m := middleware.Auth

	// Routes
	g.POST("/login", c.Login)
	g.POST("/logout", c.Logout)
	g.POST("/register", m.AuthenticateAdmin, c.Register)
}
