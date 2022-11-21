package coordinator

import (
	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/farm/counter"
	s "github.com/SimonRichardson/colossus/selectors"
)

type scanner struct {
	s.LifeCycleManager

	co      *Coordinator
	counter *counter.Farm
}

func newScanner(co *Coordinator, counter *counter.Farm) *scanner {
	return &scanner{
		LifeCycleManager: newLifeCycleService(),

		co:      co,
		counter: counter,
	}
}

func (s *scanner) Keys() ([]bs.Key, error) {
	return s.counter.Keys()
}

func (s *scanner) Size(key bs.Key) (int, error) {
	return s.counter.Size(key)
}

func (s *scanner) Members(key bs.Key) ([]bs.Key, error) {
	return s.counter.Members(key)
}
