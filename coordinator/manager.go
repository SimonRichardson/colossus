package coordinator

import (
	"github.com/SimonRichardson/colossus/coordinator/strategies"
	"github.com/SimonRichardson/colossus/farm/counter"
	"github.com/SimonRichardson/colossus/farm/notifier"
	"github.com/SimonRichardson/colossus/farm/store"
	s "github.com/SimonRichardson/colossus/selectors"
)

type manager struct {
	s.LifeCycleManager

	co       *Coordinator
	counter  *counter.Farm
	store    *store.Farm
	notifier *notifier.Farm
	strategy strategies.ManagerStrategy
	quit     chan struct{}
}

func newManager(co *Coordinator,
	cf *counter.Farm,
	sf *store.Farm,
	nf *notifier.Farm,
	strategy strategies.ManagerStrategyCreator,
) *manager {
	return &manager{
		LifeCycleManager: newLifeCycleService(),

		co:       co,
		counter:  cf,
		store:    sf,
		notifier: nf,

		strategy: strategy(co, sf),
		quit:     make(chan struct{}),
	}
}

func (m *manager) Start() error {
	go func() {
		channel := m.notifier.Subscribe(defaultInsertChannel)
		for {
			select {
			case keyFieldSize := <-channel:
				go func() {
					if err := m.strategy(keyFieldSize); err != nil {
						if err == strategies.ErrFatal {
							m.quit <- struct{}{}
						}
					}
				}()
			case <-m.quit:
				return
			}
		}
	}()
	return nil
}

func (m *manager) Stop() error {
	m.quit <- struct{}{}
	return nil
}
