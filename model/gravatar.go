package model

import (
	"fmt"
	"io"
	"crypto/md5"
)

func GetGravatar(email string, size string) string {
	u := "http://1.gravatar.com/avatar/"
	u += md5Str(email) + "?s=" + size
	return u
}

func md5Str(str string) string {
	m := md5.New()
	io.WriteString(m, str)
	return fmt.Sprintf("%x", m.Sum(nil))
}

