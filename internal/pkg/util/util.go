package util

import (
	"os"
	"path/filepath"
	"sync"

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
		"./storage",
		"./log",
	} {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, os.ModePerm)
		}
	}
}

var sliceLock sync.Mutex

func RemoveFromSlice(slice []int, target int) []int {
	sliceLock.Lock()
	var editedSlice []int
	for _, v := range slice {
		if !(v == target) {
			editedSlice = append(editedSlice, v)
		}
	}
	sliceLock.Unlock()
	return editedSlice
}

func Append(slice []int, val int) []int {
	sliceLock.Lock()
	slice = append(slice, val)
	sliceLock.Unlock()
	return slice
}

func IfExistInSlice(slice []int, target int) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}
