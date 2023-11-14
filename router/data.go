package router

import (
	"adorable-star/internal/controller"
	"adorable-star/internal/middleware"

	"github.com/gin-gonic/gin"
)

func DataRoutes(r *gin.RouterGroup) {
	g := r.Group("/data")

	// Deps
	g.Use(middleware.LongPolling)
	c := controller.Data

	// Routes
	g.POST("/fetch", c.FetchData)
	g.POST("/fetch-desc/:id", c.FetchAssignmentDesc)
	g.GET("/data", c.GetData)
	g.GET("/report", c.GetReport)
	g.GET("/course", c.GetCourses)
	g.GET("/course/:id", c.GetCourse)
	g.GET("/assignment", c.GetAssignments)
	g.GET("/assignment/:id", c.GetAssignment)
	g.GET("/message", c.GetMessages)
	g.GET("/message/:id", c.GetMessage)
	g.DELETE("/message/:id", c.DeleteMessage)
}
