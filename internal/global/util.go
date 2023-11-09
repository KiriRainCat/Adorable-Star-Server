package global

import (
	"os"
	"path/filepath"
)

func GetCwd() string {
	path, _ := os.Executable()
	return filepath.Dir(path)
}
