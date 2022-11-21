package uuid

import (
	"time"
)

func Parse(max int, duration time.Duration, perDuration int64) (*Pool, error) {
	gen, err := NewGenerator()
	if err != nil {
		return nil, err
	}
	return New(gen, max, duration, perDuration), nil
}
