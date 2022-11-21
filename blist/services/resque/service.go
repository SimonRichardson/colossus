package resque

import (
	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/typex"
)

// Tactic defines an alias for the structure of a tactic. A tactic in this
// regard consumes a slice of clusters and runs the function depending on if
// the tactic allows it.
type Tactic func([]Cluster, func(Cluster)) error

type (
	enqueueStrategy  func(*Service, Tactic) selectors.Enqueuer
	registerStrategy func(*Service, Tactic) selectors.Register
)

// EnqueueCreator defines a series of types that allow the building (creating)
// of a selector using only the service as an argument. It's expected that each
// is partially built before hand and the farm is the last argument for the
// partial application.
type (
	EnqueueCreator interface {
		Apply(*Service) selectors.Enqueuer
	}
	RegisterCreator interface {
		Apply(*Service) selectors.Register
	}
)

// DefaultService provides the bare minimum required to build the service. It
// expects that you don't need any added features and no instrumentation.
func DefaultService(address string) (*Service, error) {
	clusters, err := DefaultConfig(address)
	if err != nil {
		return nil, typex.Errorf(errors.Source, errors.InvalidArgument,
			"Invalid configuration.").With(err)
	}

	return New(clusters,
		enqueueStategyOpts{Strategy: Enqueuer, Tactic: nonBlocking},
		registerStategyOpts{Strategy: Register, Tactic: nonBlocking},
		NoopInstrumentation{},
	), nil
}

// Service defines a structure for pushing messages on to the resque message
// bus.
type Service struct {
	clusters        []Cluster
	enqueuer        selectors.Enqueuer
	register        selectors.Register
	instrumentation Instrumentation
}

// New creates a new service depending on the pool dependency
func New(clusters []Cluster,
	enq EnqueueCreator,
	reg RegisterCreator,
	instr Instrumentation,
) *Service {
	service := &Service{
		clusters:        clusters,
		instrumentation: instr,
	}

	service.enqueuer = enq.Apply(service)
	service.register = reg.Apply(service)

	return service
}

func (s *Service) EnqueueBytes(queue selectors.Queue, class selectors.Class, value []byte) error {
	return s.enqueuer.EnqueueBytes(queue, class, value)
}

func (s *Service) DequeueBytes(queue selectors.Queue, class selectors.Class) ([]byte, error) {
	return s.enqueuer.DequeueBytes(queue, class)
}

func (s *Service) RegisterFailure(queue selectors.Queue,
	class selectors.Class,
	failure selectors.Failure,
) error {
	return s.register.RegisterFailure(queue, class, failure)
}
