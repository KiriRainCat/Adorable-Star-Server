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

var lock sync.Mutex

func RemoveFromSlice(slice []int, target int) []int {
	lock.Lock()
	j := 0
	for _, v := range slice {
		if !(v == target) {
			slice[j] = v
			j++
		}
	}
	lock.Unlock()
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
