package utils

// ContainsInt64 tells whether a slice contains x.
func ContainsInt64(a []int64, x int64) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
