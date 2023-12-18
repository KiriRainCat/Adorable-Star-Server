package router

import (
	"adorable-star/internal/controller"
	"adorable-star/internal/middleware"

	"github.com/gin-gonic/gin"
)

func DataRoutes(r *gin.RouterGroup) {
	g := r.Group("/data")

	// Deps
	c := controller.Data

	// Routes
	g.POST("/fetch", middleware.Auth.AuthenticateUserLevel(1, 1), c.FetchData)
	g.POST("/fetch-desc/:id", c.FetchAssignmentDetail)

	g.GET("/report", c.GetReport)
	g.GET("/feedback-img/:id", c.GetFeedBackImage)

	g.GET("/course", c.GetCourses)
	g.GET("/course/:id", c.GetCourse)

	g.GET("/assignment", c.GetAssignments)
	g.GET("/assignment/:id", c.GetAssignment)
	g.PUT("/assignment/:id/:status", c.UpdateAssignmentStatus)
	g.POST("/assignment/upload/files/:id", c.UploadFiles)
	g.POST("/assignment/upload/juno-doc/:id", c.UploadJunoDoc)
	g.POST("/assignment/uploaded/:id", c.UnSubmit)

	g.GET("/message", c.GetMessages)
	g.GET("/message/:id", c.GetMessage)
	g.DELETE("/message", c.DeleteAllMessages)
	g.DELETE("/message/:id", c.DeleteMessage)
}
