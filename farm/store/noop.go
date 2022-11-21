package store

import (
	bs "github.com/SimonRichardson/colossus/blist/selectors"
	s "github.com/SimonRichardson/colossus/selectors"
)

// NoopSelector defines a selector that performs no operations, but attempts to
// provide "canned" data that will allow the application to still execute.
func NoopSelector(f *Farm, t Tactic) s.Selector { return noop{f} }

// NoopInserter defines a selector that performs no operations, but attempts to
// provide "canned" data that will allow the application to still execute.
func NoopInserter(f *Farm, t Tactic) s.Inserter { return noop{f} }

// NoopDeleter defines a selector that performs no operations, but attempts to
// provide "canned" data that will allow the application to still execute.
func NoopDeleter(f *Farm, t Tactic) s.Deleter { return noop{f} }

// NoopScanner defines a selector that performs no operations, but attempts to
// provide "canned" data that will allow the application to still execute.
func NoopScanner(f *Farm, t Tactic) s.Scanner { return noop{f} }

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

func (n noop) SelectRange(bs.Key, int, s.KeySizeExpiry) ([]s.KeyFieldScoreTxnValue, error) {
	return make([]s.KeyFieldScoreTxnValue, 0, 0), nil
}

func (n noop) Keys() ([]bs.Key, error) {
	return make([]bs.Key, 0, 0), nil
}

func (n noop) Size(bs.Key) (int, error) {
	return 0, nil
}

func (n noop) Members(bs.Key) ([]bs.Key, error) {
	return make([]bs.Key, 0, 0), nil
}

func (n noop) Repair([]s.KeyFieldTxnValue, s.KeySizeExpiry) error {
	return nil
}
