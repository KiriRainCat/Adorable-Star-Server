package main

import (
	"adorable-star/config"
	"adorable-star/controller"
	"adorable-star/dao"
	"adorable-star/middleware"
	"adorable-star/model"
	"adorable-star/router"
	"adorable-star/service"
	"adorable-star/service/crawler"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	//* ----------------------- Dependency Initialization ---------------------- *//
	// Init database
	db, _ := gorm.Open(sqlite.Open("./dev.db"), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	db.AutoMigrate(&model.User{}, &model.JupiterData{}, &model.Assignment{}, &model.Message{})

	// Init DAOs
	jupiterDao := dao.NewJupiterDAO(db)

	// Init middlewares
	authMiddleware := &middleware.AuthMiddleware{}

	// Init crawler
	crawler.Init(jupiterDao)

	// Create gin-engine and base router-group
	server := gin.Default()
	r := server.Group("/api", authMiddleware.Authenticate)

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
	router.AuthRoutes(r, db, authMiddleware)

	// Data APIs
	dataGroup := r.Group("/data")
	dataController := controller.NewDataController(service.NewDataService(db))
	{
		dataGroup.POST("jupiter", authMiddleware.AuthenticateAdmin, dataController.JupiterData)
	}

	server.Run(config.PORT)
}
