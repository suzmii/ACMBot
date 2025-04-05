package util

import "time"

type RateLimiter []struct {
	rate time.Duration
}
