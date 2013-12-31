package utils

import (
	"crypto/md5"
	"fmt"
)

// generate gravatar link by email.
func Gravatar(email string, size string) string {
	u := "http://1.gravatar.com/avatar/"
	u += encodeAvatarEmail(email) + "?s=" + size
	return u
}

// encode user password by sha1 with salt string from config.
func encodeAvatarEmail(email string) string {
	h := md5.New()
	h.Write([]byte(email))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
