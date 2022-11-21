package coordinator

import (
	"time"

	"github.com/SimonRichardson/colossus/blist/services/consul"
	"github.com/SimonRichardson/colossus/blist/services/consul/client"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/teleprinter"
)

type service struct {
	s.LifeCycleManager

	co        *Coordinator
	consul    *consul.Service
	frequency time.Duration
	quit      chan struct{}
}

func newService(co *Coordinator,
	consul *consul.Service,
	frequency time.Duration,
) *service {
	return &service{
		LifeCycleManager: newLifeCycleService(),

		co:        co,
		consul:    consul,
		frequency: frequency,
		quit:      make(chan struct{}),
	}
}

// Transmit pushes a value on to consul.
func (c *service) Start() error {
	errCh := make(chan error)
	go func() {
		// Don't send heartbeats!
		if c.frequency < 1 {
			errCh <- nil
			return
		}

		tick := time.NewTicker(c.frequency).C
		for {
			select {
			case <-tick:
				if err := c.consul.Heartbeat(client.Passing); err != nil {
					teleprinter.L.Error().Printf("Error sending heartbeat: %s\n", err.Error())
					continue
				}
			case <-c.quit:
				errCh <- nil
				return
			}
		}
	}()

	select {
	case err := <-errCh:
		return err
	}
}

func (c *service) Stop() error {
	c.quit <- struct{}{}
	return nil
}
