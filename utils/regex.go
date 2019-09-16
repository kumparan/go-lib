package utils

import (
	"regexp"
)

// RegexEmail returns regex for email
func RegexEmail() *regexp.Regexp {
	rgxEmail, _ := regexp.Compile("^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$")
	return rgxEmail
}
