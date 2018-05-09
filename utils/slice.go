package utils

import "strconv"

// ContainsInt64 tells whether a slice contains x.
func ContainsInt64(a []int64, x int64) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// ContainsString tells whether a slice contains x.
func ContainsString(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// SliceAtoi -> convert array of string to array of integer
func SliceAtoi(s []string) ([]int, error) {
	var arr []int

	for _, val := range s {
		i, err := strconv.Atoi(val)
		if err != nil {
			return arr, err
		}

		arr = append(arr, i)
	}
	return arr, nil
}
