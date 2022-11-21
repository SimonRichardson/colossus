package bson

import (
	"time"
)

func Parse(max int, duration time.Duration, perDuration int64) (*Pool, error) {
	return New(max, duration, perDuration), nil
}
