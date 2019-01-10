package utils

import "time"

// FormatTimeRFC3339 Format time according to RFC3339Nano
func FormatTimeRFC3339(t *time.Time) (s string) {
	if t != nil {
		s = t.Format(time.RFC3339Nano)
	}
	return
}
