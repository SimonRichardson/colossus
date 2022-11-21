package cache

import (
	"time"

	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/typex"
)

// Tactic defines an alias for the structure of a tactic. A tactic in this
// regard consumes a slice of clusters and runs the function depending on if
// the tactic allows it.
type Tactic func([]Cluster, func(Cluster)) error

type (
	encodeStrategy func(*Service, Tactic) selectors.Encoder
)

// EncodeCreator defines a series of types that allow the building (creating)
// of a selector using only the service as an argument. It's expected that each
// is partially built before hand and the farm is the last argument for the
// partial application.
type (
	EncodeCreator interface {
		Apply(*Service) selectors.Encoder
	}
)

// DefaultService provides the bare minimum required to build the service. It
// expects that you don't need any added features and no instrumentation.
func DefaultService(address string, expiry time.Duration) (*Service, error) {
	clusters, err := DefaultConfig(address, expiry)
	if err != nil {
		return nil, typex.Errorf(errors.Source, errors.InvalidArgument,
			"Invalid configuration.").With(err)
	}

	return New(clusters,
		encodeStategyOpts{Strategy: Encoder, Tactic: nonBlocking},
		NoopInstrumentation{},
	), nil
}

type Service struct {
	clusters        []Cluster
	encoder         selectors.Encoder
	instrumentation Instrumentation
}

func New(clusters []Cluster,
	enc EncodeCreator,
	instr Instrumentation,
) *Service {

	service := &Service{
		clusters:        clusters,
		instrumentation: instr,
	}

	service.encoder = enc.Apply(service)

	return service
}

func (s *Service) GetBytes(key bs.Key) ([]byte, error) {
	return s.encoder.GetBytes(key)
}

func (s *Service) SetBytes(key bs.Key, bytes []byte) error {
	return s.encoder.SetBytes(key, bytes)
}

func (s *Service) DelBytes(key bs.Key) error {
	return s.encoder.DelBytes(key)
}
