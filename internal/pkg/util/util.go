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

var appendLock sync.Mutex
var removeLock sync.Mutex

func RemoveFromSlice(slice []int, target int) []int {
	removeLock.Lock()
	j := 0
	for _, v := range slice {
		if !(v == target) {
			slice[j] = v
			j++
		}
	}
	removeLock.Unlock()
	return slice[:j]
}

func Append(slice []int, val int) []int {
	appendLock.Lock()
	slice = append(slice, val)
	appendLock.Unlock()
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
