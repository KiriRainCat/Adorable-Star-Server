package router

import (
	"adorable-star/internal/controller"
	"adorable-star/internal/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.RouterGroup) {
	g := r.Group("/admin")

	// Deps
	g.Use(middleware.Auth.AuthenticateAdmin)
	c := controller.Admin

	// Routs
	g.GET("/crawler-load", c.GetCrawlerLoad)
	g.GET("/crawler-log", c.GetCrawlerLog)
	g.GET("/sql-log", c.GetApiLog)
	g.GET("/user", c.GetUsers)
	g.PUT("/user/status/:id/:status", c.UpdateUserStatus)
	g.PUT("/user/password", c.ChangeUserPassword)
}
