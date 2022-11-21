package notifier

import (
	s "github.com/SimonRichardson/colossus/selectors"
)

// NoopNotifier defines a selector that performs no operations, but attempts to
// provide "canned" data that will allow the application to still execute.
func NoopNotifier(f *Farm, t Tactic) s.Notifier { return noop{f} }

type noop struct {
	*Farm
}

func (n noop) Publish(s.Channel, []s.KeyFieldScoreSizeExpiry) error {
	return nil
}

func (n noop) Subscribe(s.Channel) <-chan s.KeyFieldScoreSizeExpiry {
	return nil
}
