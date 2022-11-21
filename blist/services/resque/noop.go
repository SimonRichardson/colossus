package resque

import (
	"time"

	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/typex"
)

// NoopEnqueuer defines a selector that performs no operations, but attempts to
// provide "sane" data that will allow the application to still execute.
func NoopEnqueuer(s *Service, t Tactic) selectors.Enqueuer { return noop{s} }

// NoopRegister defines a selector that performs no operations, but attempts to
// provide "sane" data that will allow the application to still execute.
func NoopRegister(s *Service, t Tactic) selectors.Register { return noop{s} }

type noop struct {
	*Service
}

func (n noop) EnqueueBytes(selectors.Queue, selectors.Class, []byte) error {
	return nil
}

func (n noop) DequeueBytes(selectors.Queue, selectors.Class) ([]byte, error) {
	return nil, typex.Errorf(errors.Source, errors.MissingContent,
		"Nothing found.")
}

func (n noop) RegisterFailure(selectors.Queue,
	selectors.Class,
	selectors.Failure,
) error {
	return nil
}

type NoopInstrumentation struct{}

func (NoopInstrumentation) EnqueueCall()                          {}
func (NoopInstrumentation) EnqueueSendTo(int)                     {}
func (NoopInstrumentation) EnqueueDuration(time.Duration)         {}
func (NoopInstrumentation) DequeueCall()                          {}
func (NoopInstrumentation) DequeueSendTo(int)                     {}
func (NoopInstrumentation) DequeueDuration(time.Duration)         {}
func (NoopInstrumentation) RegisterFailureCall()                  {}
func (NoopInstrumentation) RegisterFailureSendTo(int)             {}
func (NoopInstrumentation) RegisterFailureDuration(time.Duration) {}
