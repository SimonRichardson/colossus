package persistence

import (
	p "github.com/SimonRichardson/colossus/cluster/persistence"
	"github.com/SimonRichardson/colossus/instrumentation"
	s "github.com/SimonRichardson/colossus/selectors"
)

type Tactic func([]p.Cluster, func(p.Cluster)) error

type (
	insertStrategy func(*Farm, Tactic) s.Inserter
	deleteStrategy func(*Farm, Tactic) s.Deleter
	repairStrategy func(*Farm, Tactic) s.Repairer
)

type (
	InsertCreator interface {
		Apply(*Farm) s.Inserter
	}
	DeleteCreator interface {
		Apply(*Farm) s.Deleter
	}
	RepairCreator interface {
		Apply(*Farm) s.Repairer
	}
)

// Farm defines a container for all the selectors to be able to query.
type Farm struct {
	clusters        []p.Cluster
	inserter        s.Inserter
	deleter         s.Deleter
	repairer        s.Repairer
	instrumentation instrumentation.Instrumentation
}

// New defines a function for the creation of a farm.
func New(clusters []p.Cluster,
	ins InsertCreator,
	del DeleteCreator,
	rep RepairCreator,
	instr instrumentation.Instrumentation,
) *Farm {
	farm := &Farm{
		clusters:        clusters,
		instrumentation: instr,
	}
	farm.inserter = ins.Apply(farm)
	farm.deleter = del.Apply(farm)
	farm.repairer = rep.Apply(farm)
	return farm
}

// Insert defines a way to insert some members into the store that's associated
// with the key
func (f *Farm) Insert(members []s.KeyFieldScoreTxnValue, maxSize s.KeySizeExpiry) (int, error) {
	return f.inserter.Insert(members, maxSize)
}

// Delete defines a way to delete some members into the store that's associated
// with the key
func (f *Farm) Delete(members []s.KeyFieldScoreTxnValue, maxSize s.KeySizeExpiry) (int, error) {
	return f.deleter.Delete(members, maxSize)
}

// Rollback defines a way to rollback some members into the store that's
// associated with the key
func (f *Farm) Rollback(members []s.KeyFieldScoreTxnValue, maxSize s.KeySizeExpiry) error {
	_, err := f.deleter.Delete(members, maxSize)
	return err
}

// Repair attempts to repair the store depending on the elements
func (f *Farm) Repair(elements []s.KeyFieldTxnValue, maxSize s.KeySizeExpiry) error {
	return f.repairer.Repair(elements, maxSize)
}

func (f *Farm) Topology(clusters []p.Cluster) error {
	for _, v := range f.clusters {
		if err := v.Close(); err != nil {
			return err
		}
	}

	f.clusters = clusters
	return nil
}
