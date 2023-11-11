package main

import (
	"adorable-star/internal/dao"
	"adorable-star/internal/middleware"
	"adorable-star/internal/pkg/config"
	"adorable-star/internal/pkg/crawler"
	"adorable-star/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialization
	config.Init()
	dao.Init()
	crawler.Init()
	authMiddleware := middleware.Auth

	// Create gin-engine and base router-group
	server := gin.Default()
	r := server.Group("/api")
	r.Use(authMiddleware.Authenticate).
		Use(authMiddleware.AuthenticateUser)

	//* --------------------------- API Registration --------------------------- *//
	// PING API
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "Pong",
			"data": nil,
		})
	})

	// Register API Routes
	router.AuthRoutes(r)

	server.Run(":" + config.Config.Server.Port)
}
