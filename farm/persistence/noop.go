package persistence

import (
	s "github.com/SimonRichardson/colossus/selectors"
)

// NoopInserter defines a selector that performs no operations, but attempts to
// provide "canned" data that will allow the application to still execute.
func NoopInserter(f *Farm, t Tactic) s.Inserter { return noop{f} }

// NoopDeleter defines a selector that performs no operations, but attempts to
// provide "canned" data that will allow the application to still execute.
func NoopDeleter(f *Farm, t Tactic) s.Deleter { return noop{f} }

// NoopRepairer defines a selector that performs no operations, but attempts to
// provide "canned" data that will allow the application to still execute.
func NoopRepairer(f *Farm, t Tactic) s.Repairer { return noop{f} }

type noop struct {
	*Farm
}

func (n noop) Insert([]s.KeyFieldScoreTxnValue, s.KeySizeExpiry) (int, error) {
	return 0, nil
}

func (n noop) Delete([]s.KeyFieldScoreTxnValue, s.KeySizeExpiry) (int, error) {
	return 0, nil
}

func (n noop) Rollback([]s.KeyFieldScoreTxnValue, s.KeySizeExpiry) error {
	return nil
}

func (n noop) Repair([]s.KeyFieldTxnValue, s.KeySizeExpiry) error {
	return nil
}
