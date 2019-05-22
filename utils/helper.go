package utils

// Ternary if condition is true, return a else b
func Ternary(condition bool, a, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}