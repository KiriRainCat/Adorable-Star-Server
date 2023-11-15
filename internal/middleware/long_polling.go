package middleware

import (
	"adorable-star/internal/dao"
	"adorable-star/internal/global"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func LongPolling(ctx *gin.Context) {
	// If the http method is not GET
	if ctx.Request.Method != "GET" {
		ctx.Next()
		return
	}

	// Check if uid from the token is valid
	if ctx.GetInt("uid") == 0 {
		ctx.Abort()
		return
	}

	// If header contains Instant, return data directly
	if val, err := strconv.ParseBool(ctx.Request.Header.Get("Instant")); val && err == nil {
		data, err := dao.Jupiter.GetDataByUID(ctx.GetInt("uid"))
		if err != nil {
			ctx.Abort()
			return
		}
		ctx.Set("fetchedAt", data.FetchedAt)
		ctx.Set("gpa", data.GPA)
		ctx.Next()
		return
	}

	// Wait until a value can be taken out from the channel for the user
	info := <-global.NotificationChan[ctx.GetInt("uid")]
	fetchedAt := info[0].(time.Time)
	gpa := info[1].(string)

	// Set the value got from channel to global context
	ctx.Set("fetchedAt", fetchedAt)
	ctx.Set("gpa", gpa)
	ctx.Next()
}
