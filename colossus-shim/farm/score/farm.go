package score

import (
	"time"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	c "github.com/SimonRichardson/colossus/colossus-shim/cluster/score"
	t "github.com/SimonRichardson/colossus/colossus-shim/selectors"
	"github.com/SimonRichardson/colossus/instrumentation"
)

// Tactic defines an alias for the structure of a tactic. A tactic in this
// regard consumes a slice of clusters and runs the function depending on if
// the tactic allows it.
type Tactic func([]c.Cluster, func(c.Cluster)) error

// incrementStrategy defines a way to create the selectors from a farm and a
// tactic. They have a close analogy to building (creating) something from a
// some of parts.
type (
	incrementStrategy func(*Farm, Tactic) t.Incrementer
)

// IncrementCreator defines a series of types that allow the building (creating)
// of a selector using only the farm as an argument. It's expected that each is
// partially built before hand and the farm is the last argument for the partial
// application.
type (
	IncrementCreator interface {
		Apply(*Farm) t.Incrementer
	}
)

// Farm defines a container for all the selectors to be able to query.
type Farm struct {
	clusters        []c.Cluster
	incrementer     t.Incrementer
	instrumentation instrumentation.Instrumentation
}

// New defines a function for the creation of a farm.
func New(clusters []c.Cluster,
	inc IncrementCreator,
	instr instrumentation.Instrumentation,
) *Farm {
	farm := &Farm{
		clusters:        clusters,
		instrumentation: instr,
	}
	farm.incrementer = inc.Apply(farm)
	return farm
}

// Increment defines a way to increment a score into the store that's associated
// with the key
func (f *Farm) Increment(key bs.Key, t time.Time) (int, error) {
	return f.incrementer.Increment(key, t)
}
