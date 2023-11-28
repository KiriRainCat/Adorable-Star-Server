package main

import (
	"adorable-star/internal/dao"
	"adorable-star/internal/middleware"
	"adorable-star/internal/pkg/config"
	"adorable-star/internal/pkg/crawler"
	"adorable-star/router"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialization
	config.Init()
	dao.InitDB()
	dao.InitRedis()
	crawler.Init()
	authMiddleware := middleware.Auth

	// Create gin-engine and base router-group
	gin.Default()
	server := gin.New()
	r := server.Group("/api")
	r.Use(gin.LoggerWithWriter(os.Stdout, "/api/ping")).
		Use(gin.Recovery()).
		Use(authMiddleware.Authenticate).
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
	router.DataRoutes(r)
	router.AdminRoutes(r)

	server.Run(":" + config.Config.Server.Port)
}
