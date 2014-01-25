package utils

import (
	"fmt"
	"os"
)

func FileSize(size int64) string {
	s := float64(size)
	if s > 1024*1024 {
		return fmt.Sprintf("%.1f M", s/(1024*1024))
	}
	if s > 1024 {
		return fmt.Sprintf("%.1f K", s/1024)
	}
	return fmt.Sprintf("%f B", s)
}

func IsFile(path string) bool {
	f, e := os.Stat(path)
	if e != nil {
		return false
	}
	if f.IsDir() {
		return false
	}
	return true
}

func IsDir(path string) bool {
	f, e := os.Stat(path)
	if e != nil {
		return false
	}
	return f.IsDir()
}
