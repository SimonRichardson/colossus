package permitters

import (
	"time"

	"github.com/tsenart/tb"
)

// Permitter defines if a request is allowed to succeed or not. Depending on the
// values sent to the permitter will tell you if it''' work or not.
type Permitter interface {
	Allowed(int64) bool
}

// New creates a permitter
func New(n int64, d time.Duration) Permitter {
	if n <= 0 {
		return allowAllPermitter{}
	}
	return tokenBucketPermitter{tb.NewBucket(n, d)}
}

type tokenBucketPermitter struct{ *tb.Bucket }

func (p tokenBucketPermitter) Allowed(n int64) bool {
	if value := p.Bucket.Take(n); value < n {
		p.Bucket.Put(value)
		return false
	}

	return true
}

type throttlePermitter struct {
	*tb.Bucket
	duration time.Duration
}

// NewThrottle defines a permitter that will allow all requests to be allowed,
// but instead will sleep until it's allowed to permit.
func NewThrottle(n int64, duration time.Duration) Permitter {
	if n <= 0 {
		return allowAllPermitter{}
	}
	return throttlePermitter{tb.NewBucket(n, duration), duration}
}

func (p throttlePermitter) Allowed(n int64) bool {
	if value := p.Bucket.Take(n); value < n {
		p.Bucket.Put(value)
		time.Sleep(p.duration)
	}
	return true
}

type allowAllPermitter struct{}

func (p allowAllPermitter) Allowed(n int64) bool { return true }

type allowNonePermitter struct{}

func (p allowNonePermitter) Allowed(n int64) bool { return false }
