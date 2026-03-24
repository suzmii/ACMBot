package util

import "github.com/suzmii/ACMBot/errorx/usererr"

func Zero[T any]() T {
	var zero T
	return zero
}

func ValidateRating(v int) error {
	if v < 800 || v > 3500 || v%100 != 0 {
		return usererr.ErrInvalidRating
	}
	return nil
}

func IsNumber(s string) bool {
	if s == "" {
		return false
	}
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}