package utils

import (
	"errors"
	"math/rand"
	"strconv"
	"time"
)

// PageOrDefault -> Check page value, and return 1 if page is not defined
// DEPRECATED never use this no more
func PageOrDefault(page int64) int64 {
	if page < 1 {
		return 1
	}
	return page
}

// LimitOrDefault -> Check limit value.
// DEPRECATED never use this no more
func LimitOrDefault(limit int64) (int64, error) {
	if limit == 0 {
		return 10, nil
	}
	if limit < 1 || limit > 100 {
		return 0, errors.New("Limit Value should be between 1 and 100")
	}
	return limit, nil
}

// Int642String :nodoc:
// DEPRECATED never use this no more
func Int642String(i int64) string {
	s := strconv.FormatInt(i, 10)
	return s
}

// Offset to get offset from page and limit, min value for page = 1
func Offset(page, limit int64) int64 {
	offset := (page - 1) * limit
	if offset < 0 {
		return 0
	}
	return offset
}

// GenerateID based on current time
func GenerateID() int64 {
	return int64(time.Now().UnixNano()) + int64(rand.Intn(10000))
}

// Int32PointerToInt64 :nodoc:
func Int32PointerToInt64(i *int32) int64 {
	if i != nil {
		return int64(*i)
	}
	return int64(0)
}

// Int32PointerToInt32 :nodoc:
func Int32PointerToInt32(i *int32) int32 {
	if i != nil {
		return *i
	}
	return 0
}

// IsSameSliceIgnoreOrder to compare slice without order
func IsSameSliceIgnoreOrder(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	diff := make(map[int64]bool, len(a))
	for _, v := range a {
		diff[v] = true
	}
	for _, v := range b {
		if _, ok := diff[v]; !ok {
			return false
		}
		delete(diff, v)
	}
	if len(diff) == 0 {
		return true
	}
	return false
}
