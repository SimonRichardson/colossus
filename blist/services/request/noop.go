package request

import (
	"net/http"
	"time"

	"github.com/SimonRichardson/colossus/blist/selectors"
)

// NoopRequest defines a selector that performs no operations, but attemps to
// provide "sane" data that will allow the application to still execute.
func NoopRequest(s *Service, t Tactic) selectors.Request { return noop{s} }

type noop struct {
	*Service
}

func (n noop) Request(*http.Request) (*http.Response, error) {
	return nil, nil
}

type NoopInstrumentation struct{}

func (NoopInstrumentation) RequestCall()                  {}
func (NoopInstrumentation) RequestSendTo(int)             {}
func (NoopInstrumentation) RequestDuration(time.Duration) {}
func (NoopInstrumentation) RequestRetrieved(int)          {}
func (NoopInstrumentation) RequestReturned(int)           {}
