package pool

import "fmt"

var (
	ErrTypeError   = fmt.Errorf("Type miss-match")
	ErrInvalidSeed = fmt.Errorf("Invalid Seed")
	ErrExhaustion  = fmt.Errorf("Pool Exhaustion")
	ErrRateLimited = fmt.Errorf("Rate Limited")
)
