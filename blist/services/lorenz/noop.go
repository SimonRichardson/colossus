package lorenz

import (
	"time"

	"github.com/SimonRichardson/colossus/blist/selectors"
)

// NoopCharger defines a selector that performs no operations, but attempts to
// provide "sane" data that will allow the application to still execute.
func NoopCharger(s *Service, t Tactic) selectors.Charger { return noop{s} }

// NoopEventSelector defines a selector that performs no operations, but
// attempts to provide "canned" data that will allow the application to still
// execute.
func NoopEventSelector(s *Service, t Tactic) selectors.EventSelector { return noop{s} }

// NoopCodeSetSelector defines a selector that performs no operations, but
// attempts to provide "canned" data that will allow the application to still
// execute.
func NoopCodeSetSelector(s *Service, t Tactic) selectors.CodeSetSelector { return noop{s} }

// NoopInspector defines a selector that performs no operations, but attempts to
// provide "sane" data that will allow the application to still execute.
func NoopInspector(s *Service, t Tactic) selectors.Inspector { return noop{s} }

type noop struct {
	*Service
}

func (n noop) Charge(selectors.Event,
	selectors.User,
	selectors.Payment,
) (selectors.Key, error) {
	return selectors.Key(""), nil
}

func (n noop) SelectEventByKey(selectors.Key) (selectors.Event, error) {
	return selectors.Event{}, nil
}

func (n noop) SelectEventsByOffset(int, int) ([]selectors.Event, error) {
	return nil, nil
}

func (n noop) SelectCodeForEvent(selectors.Event, selectors.User) (selectors.CodeSet, error) {
	return selectors.CodeSet{}, nil
}

func (n noop) Version() (map[string][]selectors.Version, error) {
	return map[string][]selectors.Version{}, nil
}

type NoopInstrumentation struct{}

func (NoopInstrumentation) EventSelectCall()                  {}
func (NoopInstrumentation) EventSelectSendTo(int)             {}
func (NoopInstrumentation) EventSelectDuration(time.Duration) {}
func (NoopInstrumentation) CodeSelectCall()                   {}
func (NoopInstrumentation) CodeSelectSendTo(int)              {}
func (NoopInstrumentation) CodeSelectDuration(time.Duration)  {}
func (NoopInstrumentation) InspectCall()                      {}
func (NoopInstrumentation) InspectSendTo(int)                 {}
func (NoopInstrumentation) InspectDuration(time.Duration)     {}
func (NoopInstrumentation) InspectRetrieved(int)              {}
func (NoopInstrumentation) InspectReturned(int)               {}
func (NoopInstrumentation) ChargeCall()                       {}
func (NoopInstrumentation) ChargeSendTo(int)                  {}
func (NoopInstrumentation) ChargeDuration(time.Duration)      {}
func (NoopInstrumentation) ChargeRetrieved(int)               {}
func (NoopInstrumentation) ChargeReturned(int)                {}
