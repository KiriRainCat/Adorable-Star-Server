package util

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func Log(scope string, format string, v ...any) {
	writer, _ := os.OpenFile(GetCwd()+"/storage/log/"+scope+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	defer writer.Close()

	if gin.Mode() != gin.ReleaseMode {
		log.Printf(format, v...)
	}

	log.New(writer, "\r\n", log.Ldate|log.Ltime).Printf(format, v...)
}
