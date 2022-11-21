package selectors

import (
	"time"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
)

// Incrementer defines a way to increment a series of items into the storage
type Incrementer interface {
	Increment(bs.Key, time.Time) (int, error)
}
