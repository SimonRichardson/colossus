package strategies

import "sync/atomic"

type roundRobin struct {
	ops int32
}

// NewRoundRobin defines a round robin strategy for accessing a pool item.
func NewRoundRobin() *roundRobin {
	return &roundRobin{0}
}

func (r *roundRobin) Select(key string, max int) int {
	return int(atomic.AddInt32(&r.ops, 1)) % max
}
