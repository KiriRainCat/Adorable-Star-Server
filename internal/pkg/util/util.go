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

func EnsureNessesaryDirs() {
	for _, dir := range []string{
		"./storage/img",
		"./storage/log",
	} {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, os.ModePerm)
		}
	}
}

func RemoveFromSlice(slice []int, target int) []int {
	j := 0
	for _, v := range slice {
		if !(v == target) {
			slice[j] = v
			j++
		}
	}
	return slice[:j]
}

func IfExistInSlice(slice []int, target int) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}
