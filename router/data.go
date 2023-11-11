package router

import (
	"adorable-star/internal/controller"

	"github.com/gin-gonic/gin"
)

func DataRoutes(r *gin.RouterGroup) {
	g := r.Group("/data")

	// Deps
	c := controller.Data

	// Routes
	g.GET("/courses", c.GetCourses)
	g.GET("/course", c.GetCourse)
	g.GET("/assignments", c.GetAssignments)
	g.GET("/assignment", c.GetAssignment)
	g.GET("/report", c.GetReport)
}
