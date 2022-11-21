package notifier

import (
	c "github.com/SimonRichardson/colossus/cluster/notifier"
	"github.com/SimonRichardson/colossus/instrumentation"
	s "github.com/SimonRichardson/colossus/selectors"
)

// Tactic defines an alias for the structure of a tactic. A tactic in this
// regard consumes a slice of clusters and runs the function depending on if
// the tactic allows it.
type Tactic func([]c.Cluster, func(c.Cluster)) error

// notifierStrategy defines a way to create the selectors from a farm and a
// tactic. They have a close analogy to building (creating) something from a
// some of parts.
type (
	notifierStrategy func(*Farm, Tactic) s.Notifier
)

// NotifyCreator defines a series of types that allow the building (creating) of
// a selector using only the farm as an argument. It's expected that each is
// partially built before hand and the farm is the last argument for the partial
// application.
type (
	NotifyCreator interface {
		Apply(*Farm) s.Notifier
	}
)

// Farm defines a container for all the selectors to be able to query.
type Farm struct {
	clusters        []c.Cluster
	notifier        s.Notifier
	instrumentation instrumentation.Instrumentation
}

// New defines a function for the creation of a farm.
func New(clusters []c.Cluster,
	not NotifyCreator,
	instr instrumentation.Instrumentation,
) *Farm {
	farm := &Farm{
		clusters:        clusters,
		instrumentation: instr,
	}
	farm.notifier = not.Apply(farm)
	return farm
}

// Publish defines a way to publish some changes that has occured recently.
func (f *Farm) Publish(channel s.Channel, members []s.KeyFieldScoreSizeExpiry) error {
	return f.notifier.Publish(channel, members)
}

// Unpublish defines a way to publish some changes that has occured recently.
func (f *Farm) Unpublish(channel s.Channel, members []s.KeyFieldScoreSizeExpiry) error {
	return f.notifier.Unpublish(channel, members)
}

// Subscribe defines a way to recieve notifications that something has been
// published to the system.
func (f *Farm) Subscribe(channel s.Channel) <-chan s.KeyFieldScoreSizeExpiry {
	return f.notifier.Subscribe(channel)
}

func (f *Farm) Topology(clusters []c.Cluster) error {
	for _, v := range f.clusters {
		if err := v.Close(); err != nil {
			return err
		}
	}

	f.clusters = clusters
	return nil
}
