package request

import (
	"net/http"

	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/typex"
)

// Tactic defines an alias for the structure of a tactic. A tactic in this
// regard consumes a slice of clusters and runs the function depending on if
// the tactic allows it.
type Tactic func([]Cluster, func(int, Cluster)) error

type (
	requestStrategy func(*Service, Tactic) selectors.Request
)

// RequestCreator defines a series of types that allow the building (creating)
// of a selector using only the service as an argument. It's expected that each
// is partially built before hand and the farm is the last argument for the
// partial application.
type (
	RequestCreator interface {
		Apply(*Service) selectors.Request
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
		First(),
		requestStategyOpts{Strategy: Request, Tactic: nonBlocking},
		NoopInstrumentation{},
	), nil
}

type Service struct {
	clusters        []Cluster
	fitting         Provider
	request         selectors.Request
	instrumentation Instrumentation
}

func New(clusters []Cluster,
	fitting Provider,
	req RequestCreator,
	instr Instrumentation,
) *Service {
	service := &Service{
		clusters:        clusters,
		fitting:         fitting,
		instrumentation: instr,
	}

	service.request = req.Apply(service)

	return service
}

func (s *Service) Request(req *http.Request) (*http.Response, error) {
	return s.request.Request(req)
}
