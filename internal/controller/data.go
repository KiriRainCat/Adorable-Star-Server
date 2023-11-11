package controller

import (
	"adorable-star/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var Data = &DataController{service.Data}

type DataController struct {
	s *service.DataService
}

func (c *DataController) GetCourses(ctx *gin.Context) {
	// Get courses
	courses, err := c.s.GetCourses(ctx.GetInt("uid"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "服务器内部发生错误，请联系开发者",
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": courses,
	})
}

func (c *DataController) GetCourse(ctx *gin.Context) {
	// Get query and check if it's empty
	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数不得为空",
			"data": nil,
		})
		return
	}

	// Get courses
	course, err := c.s.GetCourse(ctx.GetInt("uid"), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "服务器内部发生错误，请联系开发者",
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": course,
	})
}

func (c *DataController) GetAssignments(ctx *gin.Context) {
	// Get courses
	assignments, err := c.s.GetAssignments(ctx.GetInt("uid"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "服务器内部发生错误，请联系开发者",
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": assignments,
	})
}

func (c *DataController) GetAssignment(ctx *gin.Context) {
	// Get query and check if it's empty
	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数不得为空",
			"data": nil,
		})
		return
	}

	// Get courses
	assignment, err := c.s.GetAssignment(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "服务器内部发生错误，请联系开发者",
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": assignment,
	})
}

func (c *DataController) GetReport(ctx *gin.Context) {
	// Get courses
	file, err := c.s.GetReport(ctx.GetInt("uid"))
	if err != nil {
		if err.Error() == "fileNotExist" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "文件不存在，请手动检索一下数据",
				"data": nil,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "服务器内部发生错误，请联系开发者",
			"data": nil,
		})
		return
	}

	ctx.Writer.WriteString(string(file))
}

func (c *DataController) GetMessages(ctx *gin.Context) {
	// Get courses
	messages, err := c.s.GetMessages(ctx.GetInt("uid"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "服务器内部发生错误，请联系开发者",
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": messages,
	})
}

func (c *DataController) GetMessage(ctx *gin.Context) {
	// Get query and check if it's empty
	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数不得为空",
			"data": nil,
		})
		return
	}

	// Get courses
	message, err := c.s.GetMessage(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "服务器内部发生错误，请联系开发者",
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": message,
	})
}

func (c *DataController) DeleteMessage(ctx *gin.Context) {
	// Get query and check if it's empty
	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数不得为空",
			"data": nil,
		})
		return
	}

	// Get courses
	err := c.s.DeleteMessage(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "服务器内部发生错误，请联系开发者",
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": nil,
	})
}
