package utils

import (
	"regexp"
	"strings"
)

var (
	regexEmail   *regexp.Regexp
	regexASCII   *regexp.Regexp
	regexEnglish *regexp.Regexp
	regexURL     *regexp.Regexp
)

func init() {
	regexEmail, _ = regexp.Compile(`(?i)[A-Z0-9._%+-]+@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}`)
	regexASCII, _ = regexp.Compile(`^[a-zA-Z0-9-]+$`)
	regexEnglish, _ = regexp.Compile(`^[a-zA-Z]+$`)
	regexURL, _ = regexp.Compile(`(ftp|http|https):\/\/(\w+:{0,1}\w*@)?(\S+)(:[0-9]+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?`)
}

func IsEmptyString(str string) bool {
	return len(str) == 0
}

func IsEmail(str string) bool {
	return regexEmail.MatchString(str)
}

func IsURL(str string) bool {
	return regexURL.MatchString(str)
}

func IsLonger(str string, length int) bool {
	return len(str) > length
}

func IsShorter(str string, length int) bool {
	return len(str) < length
}

func IsASCII(str string) bool {
	return regexASCII.MatchString(str)
}

func IsEnglish(str string) bool {
	return regexEnglish.MatchString(str)
}

func IsContain(str string, contain string) bool {
	return strings.Contains(str, contain)
}
