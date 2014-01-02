package utils

import (
	"crypto/sha1"
	"fmt"
	"io"
)

func Sha1(raw string) string {
	t := sha1.New()
	io.WriteString(t, raw)
	return fmt.Sprintf("%x", t.Sum(nil))
}
