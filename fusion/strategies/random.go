package strategies

import "math/rand"

type random struct{}

func NewRandom() random {
	return random{}
}

func (r random) Select(key string, max int) int {
	return rand.Intn(max)
}
