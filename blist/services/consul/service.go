package consul

import (
	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	fs "github.com/SimonRichardson/colossus/fusion/selectors"
	"github.com/SimonRichardson/colossus/typex"
)

// Tactic defines an alias for the structure of a tactic. A tactic in this
// regard consumes a slice of clusters and runs the function depending on if
// the tactic allows it.
type Tactic func([]Cluster, func(Cluster)) error

type (
	semaphoreStrategy func(*Service, Tactic) selectors.Semaphore
	heartbeatStrategy func(*Service, Tactic) selectors.Heartbeat
	keyStoreStrategy  func(*Service, Tactic) selectors.KeyStore
)

// SemaphoreCreator defines a series of types that allow the building (creating)
// of a selector using only the service as an argument. It's expected that each
// is partially built before hand and the farm is the last argument for the
// partial application.
type (
	SemaphoreCreator interface {
		Apply(*Service) selectors.Semaphore
	}
	HeartbeatCreator interface {
		Apply(*Service) selectors.Heartbeat
	}
	KeyStoreCreator interface {
		Apply(*Service) selectors.KeyStore
	}
)

// DefaultService provides the bare minimum required to build the service. It
// expects that you don't need any added features and no instrumentation.
func DefaultService(address, checkId, output string) (*Service, error) {
	clusters, err := DefaultConfig(address, checkId, output)
	if err != nil {
		return nil, typex.Errorf(errors.Source, errors.InvalidArgument,
			"Invalid configuration.").With(err)
	}

	return New(clusters,
		semaphoreStategyOpts{Strategy: NoopSemaphore, Tactic: nonBlocking},
		heartbeatStategyOpts{Strategy: Heartbeat, Tactic: nonBlocking},
		keyStoreStategyOpts{Strategy: NoopKeyStore, Tactic: nonBlocking},
		NoopInstrumentation{},
	), nil
}

type Service struct {
	clusters        []Cluster
	semaphore       selectors.Semaphore
	heartbeat       selectors.Heartbeat
	keyStore        selectors.KeyStore
	instrumentation Instrumentation
}

func New(clusters []Cluster,
	sem SemaphoreCreator,
	hrt HeartbeatCreator,
	kvs KeyStoreCreator,
	instr Instrumentation,
) *Service {
	service := &Service{
		clusters:        clusters,
		instrumentation: instr,
	}

	service.semaphore = sem.Apply(service)
	service.heartbeat = hrt.Apply(service)
	service.keyStore = kvs.Apply(service)

	return service
}

func (s *Service) Lock(ns selectors.Namespace) (selectors.SemaphoreUnlock, error) {
	return s.semaphore.Lock(ns)
}

func (s *Service) Heartbeat(status selectors.HealthStatus) error {
	return s.heartbeat.Heartbeat(status)
}

func (s *Service) List(prefix fs.Prefix) (map[string]int, error) {
	return s.keyStore.List(prefix)
}
