package utils

import "strings"

// StandardizeSpaces -> Join long query to one line query
func StandardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
