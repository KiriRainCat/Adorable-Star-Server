package util

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func GetCwd() string {
	var path string
	if gin.Mode() == gin.ReleaseMode {
		path, _ = os.Executable()
	}
	return filepath.Dir(path)
}
