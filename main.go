package main

import (
	"adorable-star/config"
	"adorable-star/controller"
	"adorable-star/crawler"
	"adorable-star/middleware"
	"adorable-star/model"
	"adorable-star/router"
	"adorable-star/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	// Init necessary deps
	crawler.Init()
	db, _ := gorm.Open(sqlite.Open("./dev.db"), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	db.AutoMigrate(&model.User{}, &model.Assignment{}, &model.Message{})

	// Create gin-engine and base router-group
	server := gin.Default()
	r := server.Group("/api", (&middleware.AppAuthMiddleware{}).Authenticate)

	// PING API
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "Pong",
			"data": nil,
		})
	})

	// Register API Routes
	router.AuthRoutes(r, db)

	// Data APIs
	dataGroup := r.Group("/data")
	dataController := controller.NewDataController(service.NewDataService(db))
	{
		dataGroup.POST("jupiter", dataController.JupiterData)
	}

	server.Run(config.PORT)
}
