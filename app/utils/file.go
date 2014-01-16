package utils

import (
	"fmt"
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
