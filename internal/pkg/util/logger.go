package util

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func Log(scope string, format string, v ...any) {
	writer, _ := GetFileWriter("storage/log/" + scope + ".log")
	defer writer.Close()

	if gin.Mode() != gin.ReleaseMode {
		log.Printf(format, v...)
	}

	log.New(writer, "", log.Ldate|log.Ltime).Printf(format, v...)
}

// @Param `path` : Relative path without leading '/'
func GetFileWriter(path string) (*os.File, error) {
	return os.OpenFile(GetCwd()+"/"+path, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
}
