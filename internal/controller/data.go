package controller

import (
	"adorable-star/internal/pkg/response"
	"adorable-star/internal/pkg/util"
	"adorable-star/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var Data = &DataController{service.Data}

// Rate limiter lists
var fetchDataRateLimiter = []int{}
var fetchDetailRateLimiter = []int{}
var UploadRateLimiter = []int{}

type DataController struct {
	s *service.DataService
}

func (c *DataController) FetchData(ctx *gin.Context) {
	// Limit fetch rate
	uid := ctx.GetInt("uid")
	if !util.IfExistInSlice(fetchDataRateLimiter, uid) {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"code": http.StatusTooManyRequests,
			"msg":  "请求过于频繁，请稍后再试",
			"data": nil,
		})
		return
	}
	fetchDataRateLimiter = append(fetchDataRateLimiter, uid)
	defer util.RemoveFromSlice(fetchDataRateLimiter, uid)

	// Fetch data
	c.s.FetchData(uid)

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": nil,
	})
}

func (c *DataController) FetchAssignmentDetail(ctx *gin.Context) {
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

	// Limit fetch rate
	uid := ctx.GetInt("uid")
	if !util.IfExistInSlice(fetchDetailRateLimiter, uid) {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"code": http.StatusTooManyRequests,
			"msg":  "请求过于频繁，请稍后再试",
			"data": nil,
		})
		return
	}
	fetchDetailRateLimiter = append(fetchDetailRateLimiter, uid)
	defer util.RemoveFromSlice(fetchDetailRateLimiter, uid)

	// Fetch assignment detail
	force, _ := strconv.ParseBool(ctx.Query("force"))
	err := c.s.FetchAssignmentDetail(uid, id, force)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "服务器内部发生错误，请联系开发者",
			"data": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": nil,
	})
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
		"data": response.Data{
			FetchedAt: ctx.GetTime("fetchedAt"),
			GPA:       ctx.GetString("gpa"),
			Data:      courses,
		},
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

	// Get course
	course, err := c.s.GetCourse(ctx.GetInt("uid"), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "参数错误",
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

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": response.Data{
			FetchedAt: ctx.GetTime("fetchedAt"),
			GPA:       ctx.GetString("gpa"),
			Data:      course,
		},
	})
}

func (c *DataController) GetAssignments(ctx *gin.Context) {
	// Get assignments
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
		"data": response.Data{
			FetchedAt: ctx.GetTime("fetchedAt"),
			GPA:       ctx.GetString("gpa"),
			Data:      assignments,
		},
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

	// Get assignment
	assignment, err := c.s.GetAssignment(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "参数错误",
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

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": response.Data{
			FetchedAt: ctx.GetTime("fetchedAt"),
			GPA:       ctx.GetString("gpa"),
			Data:      assignment,
		},
	})
}

func (c *DataController) GetReport(ctx *gin.Context) {
	// Get report card
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
	// Get messages
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
		"data": response.Data{
			FetchedAt: ctx.GetTime("fetchedAt"),
			GPA:       ctx.GetString("gpa"),
			Data:      messages,
		},
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

	// Get message
	message, err := c.s.GetMessage(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "参数错误",
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

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": message,
	})
}

func (c *DataController) UpdateAssignmentStatus(ctx *gin.Context) {
	// Get query and check if it's empty
	id, _ := strconv.Atoi(ctx.Param("id"))
	status, _ := strconv.Atoi(ctx.Param("status"))
	if id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数不得为空",
			"data": nil,
		})
		return
	}

	// Get assignment
	err := c.s.UpdateAssignmentStatus(id, status)
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
		"data": response.Data{
			FetchedAt: ctx.GetTime("fetchedAt"),
			GPA:       ctx.GetString("gpa"),
			Data:      nil,
		},
	})
}

func (c *DataController) UploadJunoDoc(ctx *gin.Context) {
	type json struct {
		Text string `json:"text" binding:"required"`
	}

	// Get query and check if it's empty
	var data *json
	if ctx.BindJSON(&data) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}

	// Limit upload rate
	uid := ctx.GetInt("uid")
	if !util.IfExistInSlice(UploadRateLimiter, uid) {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"code": http.StatusTooManyRequests,
			"msg":  "请求过于频繁，请稍后再试",
			"data": nil,
		})
		return
	}
	UploadRateLimiter = append(UploadRateLimiter, uid)
	defer util.RemoveFromSlice(UploadRateLimiter, uid)

	// Turn in JunoDoc to Jupiter Ed
	if err := c.s.TurnInJunoDoc(uid, id, data.Text); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "提交失败",
			"data": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": nil,
	})
}

func (c *DataController) UploadFiles(ctx *gin.Context) {
	// Get query and check if it's empty
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "文件异常",
			"data": nil,
		})
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}

	// Limit upload rate
	uid := ctx.GetInt("uid")
	if !util.IfExistInSlice(UploadRateLimiter, uid) {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"code": http.StatusTooManyRequests,
			"msg":  "请求过于频繁，请稍后再试",
			"data": nil,
		})
		return
	}
	UploadRateLimiter = append(UploadRateLimiter, uid)
	defer util.RemoveFromSlice(UploadRateLimiter, uid)

	// Save files
	for _, file := range form.File["files"] {
		ctx.SaveUploadedFile(file, util.GetCwd()+"/storage/tmp"+strconv.Itoa(uid)+"/"+file.Filename)
	}

	// Turn in files to Jupiter Ed
	if err := c.s.TurnInFiles(uid, id); err != nil {
		println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "提交失败",
			"data": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": nil,
	})
}

func (c *DataController) UnSubmit(ctx *gin.Context) {
	type json struct {
		Name string `json:"name" binding:"required"`
	}

	// Get query and check if it's empty
	var data *json
	id, _ := strconv.Atoi(ctx.Param("id"))
	if ctx.BindJSON(&data) != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}

	//
	err := c.s.UnSubmit(ctx.GetInt("uid"), id, data.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "服务器内部发生错误，请联系开发者",
			"data": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": nil,
	})
}

func (c *DataController) DeleteAllMessages(ctx *gin.Context) {
	// Delete messages
	err := c.s.DeleteAllMessages(ctx.GetInt("uid"))
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

	// Delete message
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
