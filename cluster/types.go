package cluster

import (
	bs "github.com/SimonRichardson/colossus/blist/selectors"
	s "github.com/SimonRichardson/colossus/selectors"
)

// Inserter represents a way to insert a mass collection of members in to the
// store. This is slightly different setup to the selectors interface to enable
// better concurrency.
type Inserter interface {
	Insert([]s.KeyFieldScoreTxnValue, s.KeySizeExpiry) <-chan Element
}

// Modifier represents a way to insert a mass collection of members in to the
// store. This is slightly different setup to the selectors interface to enable
// better concurrency.
type Modifier interface {
	Modify([]s.KeyFieldScoreTxnValue, s.KeySizeExpiry) <-chan Element
}

// Deleter represents a way to delete a mass collection of members in to the
// store. This is slightly different setup to the selectors interface to enable
// better concurrency.
type Deleter interface {
	Delete([]s.KeyFieldScoreTxnValue, s.KeySizeExpiry) <-chan Element
	Rollback([]s.KeyFieldScoreTxnValue, s.KeySizeExpiry) <-chan Element
}

// Selector defines a way to select members from the store.
type Selector interface {
	Select(bs.Key, bs.Key) <-chan Element
	SelectRange(bs.Key, int, s.KeySizeExpiry) <-chan Element
}

// Scanner represents a way to introspect the store to help understand what the
// store has with in it's collections.
type Scanner interface {
	Keys() <-chan Element
	Size(bs.Key) <-chan Element
	Members(bs.Key) <-chan Element
}

// Scorer defines a way to score members with in the collection.
type Scorer interface {
	Score([]s.KeyFieldTxnValue) (map[s.KeyFieldTxnValue]s.Presence, error)
}

// Repairer defines a way to *attempt* to repair the collection, if possible.
type Repairer interface {
	Repair([]s.KeyFieldScoreTxnValue, s.KeySizeExpiry) <-chan Element
}

// Notifier defines a way to publish various messages to a channel
type Notifier interface {
	Publish(s.Channel, []s.KeyFieldScoreSizeExpiry) <-chan Element
	Unpublish(s.Channel, []s.KeyFieldScoreSizeExpiry) <-chan Element
	Subscribe(s.Channel) <-chan Element
}

// Closer closes the current cluster along with any underlying pools.
type Closer interface {
	Close() error
}
