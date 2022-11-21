package pool

import (
	"time"

	"github.com/tsenart/tb"
)

type Tactic func(int, func() error) error

func nonBlocking() Tactic {
	return func(amount int, fn func() error) error {
		return fn()
	}
}

func rateLimited(perDuration int64, duration time.Duration) Tactic {
	permits := Permitter{tb.NewBucket(perDuration, duration)}
	return func(amount int, fn func() error) error {
		if permits.Allowed(int64(amount)) {
			return fn()
		}
		return ErrRateLimited
	}
}

func GetTactic(duration time.Duration, perDuration int64) Tactic {
	if perDuration < 1 {
		return nonBlocking()
	}

	return rateLimited(perDuration, duration)
}

type Permitter struct {
	*tb.Bucket
}

func (p Permitter) Allowed(n int64) bool {
	if value := p.Bucket.Take(n); value < n {
		p.Bucket.Put(value)
		return false
	}

	return true
}
