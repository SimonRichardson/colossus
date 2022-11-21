package lorenz

import (
	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/typex"
)

const (
	ServiceName = "stripe"

	defaultServiceName = "lorenz"
)

// Tactic defines an alias for the structure of a tactic. A tactic in this
// regard consumes a slice of clusters and runs the function depending on if
// the tactic allows it.
type Tactic func([]Cluster, func(Cluster)) error

type (
	chargeStrategy          func(*Service, Tactic) selectors.Charger
	eventSelectorStrategy   func(*Service, Tactic) selectors.EventSelector
	codeSetSelectorStrategy func(*Service, Tactic) selectors.CodeSetSelector
	inspectStrategy         func(*Service, Tactic) selectors.Inspector
)

// InspectCreator defines a series of types that allow the building (creating)
// of a selector using only the service as an argument. It's expected that each
// is partially built before hand and the farm is the last argument for the
// partial application.
type (
	ChargeCreator interface {
		Apply(*Service) selectors.Charger
	}
	EventSelectorCreator interface {
		Apply(*Service) selectors.EventSelector
	}
	CodeSetSelectorCreator interface {
		Apply(*Service) selectors.CodeSetSelector
	}
	InspectCreator interface {
		Apply(*Service) selectors.Inspector
	}
)

// DefaultService provides the bare minimum required to build the service. It
// expects that you don't need any added features and no instrumentation.
func DefaultService(addresses, version string) (*Service, error) {
	clusters, err := DefaultConfig(addresses, version)
	if err != nil {
		return nil, typex.Errorf(errors.Source, errors.InvalidArgument,
			"Invalid configuration.").With(err)
	}

	return New(clusters,
		chargeStategyOpts{Strategy: Charger, Tactic: nonBlocking},
		eventSelectStategyOpts{Strategy: EventSelector, Tactic: nonBlocking},
		codeSetSelectStategyOpts{Strategy: CodeSetSelector, Tactic: nonBlocking},
		inspectStategyOpts{Strategy: Inspector, Tactic: nonBlocking},
		NoopInstrumentation{},
	), nil
}

type Service struct {
	clusters        []Cluster
	charger         selectors.Charger
	events          selectors.EventSelector
	codes           selectors.CodeSetSelector
	inspector       selectors.Inspector
	instrumentation Instrumentation
}

func New(clusters []Cluster,
	cha ChargeCreator,
	evt EventSelectorCreator,
	cod CodeSetSelectorCreator,
	ins InspectCreator,
	instr Instrumentation,
) *Service {

	service := &Service{
		clusters:        clusters,
		instrumentation: instr,
	}

	service.charger = cha.Apply(service)
	service.events = evt.Apply(service)
	service.codes = cod.Apply(service)
	service.inspector = ins.Apply(service)

	return service
}

func (s *Service) Charge(event selectors.Event,
	user selectors.User,
	element selectors.Payment,
) (selectors.Key, error) {
	return s.charger.Charge(event, user, element)
}

func (s *Service) SelectEventByKey(key selectors.Key) (selectors.Event, error) {
	return s.events.SelectEventByKey(key)
}

func (s *Service) SelectEventsByOffset(offset, limit int) ([]selectors.Event, error) {
	return s.events.SelectEventsByOffset(offset, limit)
}

func (s *Service) SelectCodeForEvent(event selectors.Event, user selectors.User) (selectors.CodeSet, error) {
	return s.codes.SelectCodeForEvent(event, user)
}

func (s *Service) Version() (map[string][]selectors.Version, error) {
	return s.inspector.Version()
}
