package utils

import (
	"time"

	"github.com/jpillora/backoff"
)

type stop struct {
	error
}

var backoffer = &backoff.Backoff{
	Min:    200 * time.Millisecond,
	Max:    1 * time.Second,
	Jitter: true,
}

// Retry :nodoc:
func Retry(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}

		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return Retry(attempts, 2*sleep, fn)
		}
		return err
	}
	return nil
}
