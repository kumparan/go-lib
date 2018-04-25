package utils

import "strings"

// StandardizeSpaces -> Join long query to one line query
func StandardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// GenerateSlug  -> Replace space with dash, lower it, and trim the space
func GenerateSlug(inputStr string) string {
	return strings.Trim(strings.ToLower(strings.Replace(inputStr, " ", "-", -1)), " ")
}
