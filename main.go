package main

import (
	"adorable-star/config"
	"adorable-star/internal/controller"
	"adorable-star/internal/dao"
	"adorable-star/internal/middleware"
	"adorable-star/internal/pkg/crawler"
	"adorable-star/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Init database
	dao.Init()

	// Launch Crawler
	crawler.Init()
	authMiddleware := middleware.Auth

	// Create gin-engine and base router-group
	server := gin.Default()
	r := server.Group("/api", authMiddleware.Authenticate, authMiddleware.AuthenticateUser)

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

	// Data APIs
	dataGroup := r.Group("/data")
	dataController := controller.Data
	{
		dataGroup.POST("jupiter", authMiddleware.AuthenticateAdmin, dataController.JupiterData)
	}

	server.Run(config.PORT)
}
