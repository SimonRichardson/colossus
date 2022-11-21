package request

import (
	"net/http"
	"sync"

	sv "github.com/SimonRichardson/colossus/blist/services"
	"github.com/SimonRichardson/colossus/blist/services/request/client"
)

// Cluster defines a interface for parallelised requesting.
type Cluster interface {
	sv.Request
	sv.Index
}

type cluster struct {
	pools *client.Pools
	index int
}

func newCluster(pools *client.Pools, index int) *cluster {
	return &cluster{pools, index}
}

func (c *cluster) Request(req *http.Request) <-chan sv.Element {
	return c.common(func(cli client.Client, dst chan sv.Element) {
		if resp, err := cli.Request(req); err != nil {
			dst <- sv.NewErrorElement(err)
		} else {
			dst <- sv.NewResponseElement(c.index, resp)
		}
	})
}

func (c *cluster) Index() int {
	return c.index
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
