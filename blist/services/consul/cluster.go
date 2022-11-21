package consul

import (
	"sync"

	"github.com/SimonRichardson/colossus/blist/selectors"
	sv "github.com/SimonRichardson/colossus/blist/services"
	"github.com/SimonRichardson/colossus/blist/services/consul/client"
	fs "github.com/SimonRichardson/colossus/fusion/selectors"
)

// Cluster defines a interface for parallelised requesting.
type Cluster interface {
	sv.Semaphore
	sv.Heartbeat
	sv.KeyStore
}

type cluster struct {
	pools *client.Pools
}

func newCluster(pools *client.Pools) *cluster {
	return &cluster{pools}
}

func (c *cluster) Lock(ns selectors.Namespace) <-chan sv.Element {
	return c.common(func(cli client.Client, dst chan sv.Element) {
		if fn, err := lock(cli, ns); err != nil {
			dst <- sv.NewErrorElement(err)
		} else {
			dst <- sv.NewSemaphoreUnlockElement(fn)
		}
	})
}

func (c *cluster) Heartbeat(h selectors.HealthStatus) <-chan sv.Element {
	return c.common(func(cli client.Client, dst chan sv.Element) {
		err := beat(cli, h)
		dst <- sv.NewErrorElement(err)
	})
}

func (c *cluster) List(ns fs.Prefix) <-chan sv.Element {
	return c.common(func(cli client.Client, dst chan sv.Element) {
		if fn, err := list(cli, ns); err != nil {
			dst <- sv.NewErrorElement(err)
		} else {
			dst <- sv.NewMapStringIntElement(fn)
		}
	})
}

func (c *cluster) common(fn func(client.Client, chan sv.Element)) <-chan sv.Element {
	out := make(chan sv.Element)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() { wg.Wait(); close(out) }()
		go func() {
			defer wg.Done()

			if err := c.pools.With(func(cli client.Client) error {
				fn(cli, out)
				return nil
			}); err != nil {
				out <- sv.NewErrorElement(err)
			}
		}()
	}()
	return out
}

func lock(cli client.Client, ns selectors.Namespace) (selectors.SemaphoreUnlock, error) {
	return cli.Lock(ns)
}

func beat(cli client.Client, h selectors.HealthStatus) error {
	return cli.Heartbeat(h)
}

func list(cli client.Client, p fs.Prefix) (map[string]int, error) {
	return cli.List(p)
}
