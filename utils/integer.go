package utils

import "errors"

// PageOrDefault -> Check page value, and return 1 if page is not defined
func PageOrDefault(page int64) int64 {
	if page < 1 {
		return 1
	}
	return page
}

// LimitOrDefault -> Check limit value.
func LimitOrDefault(limit int64) (int64, error) {
	if limit == 0 {
		return 10, nil
	}
	if limit < 1 || limit > 100 {
		return 0, errors.New("Limit Value should be between 1 and 100")
	}
	return limit, nil
}
