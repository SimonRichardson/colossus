package cluster

import (
	bs "github.com/SimonRichardson/colossus/blist/selectors"
	c "github.com/SimonRichardson/colossus/cluster"
)

// Incrementer defines a way to increment a value with in the store
type Incrementer interface {
	Increment(bs.Key) <-chan c.Element
}
