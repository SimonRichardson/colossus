package cache

import (
	"time"

	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/typex"
)

// NoopEncoder defines a selector that performs no operations, but attempts to
// provide "sane" data that will allow the application to still execute.
func NoopEncoder(s *Service, t Tactic) selectors.Encoder { return noop{s} }

type noop struct {
	*Service
}

func (n noop) GetBytes(bs.Key) ([]byte, error) {
	return nil, typex.Errorf(errors.Source, errors.MissingContent,
		"Not found.")
}

func (n noop) SetBytes(bs.Key, []byte) error {
	return nil
}

func (n noop) DelBytes(bs.Key) error {
	return nil
}

type NoopInstrumentation struct{}

func (NoopInstrumentation) EncodeCall()                  {}
func (NoopInstrumentation) EncodeSendTo(int)             {}
func (NoopInstrumentation) EncodeDuration(time.Duration) {}
