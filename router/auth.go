package router

import (
	"adorable-star/controller"
	"adorable-star/middleware"
	"adorable-star/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoutes(r *gin.RouterGroup, db *gorm.DB, m *middleware.AuthMiddleware) {
	g := r.Group("/auth")

	// Deps
	authController := controller.NewAuthController(service.NewUserService(db))

	// Routes
	g.POST("/login", authController.Login)
	g.POST("/logout", authController.Logout)
	g.POST("/register", m.AuthenticateAdmin, authController.Register)
}
