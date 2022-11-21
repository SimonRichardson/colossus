package selectors

import (
	s "github.com/SimonRichardson/colossus/blist/selectors"
)

// Inserter defines a way to insert a series of items into the storage
type Inserter interface {
	Insert([]KeyFieldScoreTxnValue, KeySizeExpiry) (int, error)
}

// Modifier defines a way to modify values already existing with in the storage
// system. Essentially this boils down to a new insert that over-writes existing
// values.
type Modifier interface {
	Modify([]KeyFieldScoreTxnValue, KeySizeExpiry) (int, error)
	ModifyWithOperations(s.Key, s.Key, []Operation, float64, SizeExpiry) (int, error)
}

// Deleter defines a way to remove items that where set with in the storage
type Deleter interface {
	Delete([]KeyFieldScoreTxnValue, KeySizeExpiry) (int, error)
	Rollback([]KeyFieldScoreTxnValue, KeySizeExpiry) error
}

// Selector defines a way to query the storage
type Selector interface {
	Select(s.Key, s.Key) (KeyFieldScoreTxnValue, error)
	SelectRange(s.Key, int, KeySizeExpiry) ([]KeyFieldScoreTxnValue, error)
	// TODO (Implement SelectOffset)
}

// Scanner defines a way to introspect the storage
type Scanner interface {
	Keys() ([]s.Key, error)
	Size(s.Key) (int, error)
	Members(s.Key) ([]s.Key, error)
}

// Repairer defines a way to repair the storage
// This is mainly for internal use, but could be used as a peridoical repairing
// stragegy
type Repairer interface {
	Repair([]KeyFieldTxnValue, KeySizeExpiry) error
}

// Notifier defines a way to know when a change in the system has occured
type Notifier interface {
	Publish(Channel, []KeyFieldScoreSizeExpiry) error
	Unpublish(Channel, []KeyFieldScoreSizeExpiry) error
	Subscribe(Channel) <-chan KeyFieldScoreSizeExpiry
}

// Manager defines a way to start something then stop something with in a system
type Manager interface {
	Start() error
	Stop() error
}

// Inspector defines a way to inspect the store.
// Note: it's not optimised and can be considered exploitative and slow
type Inspector interface {
	Query(s.Key, QueryOptions, SizeExpiry) ([]QueryRecord, error)
}

// LifeCycleManager defines a way to know approximately who's calling what, so
// it's then possible to know if it's safe to shut down the manager.
type LifeCycleManager interface {
	In()
	Out()
	Quit() error
}
