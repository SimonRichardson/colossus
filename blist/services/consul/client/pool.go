package client

import (
	"sync"
)

type Pools struct {
	pool *Pool
}

func New(address, checkId, output string,
	maxPoolSize int,
	creator ClientCreator,
) *Pools {
	return &Pools{
		pool: newPool(func() Client {
			return creator(address, checkId, output)
		}, maxPoolSize),
	}
}

func (p Pools) With(fn func(Client) error) error {
	client, err := p.pool.Get()
	defer p.pool.Put(client)

	if err != nil {
		return err
	}

	return fn(client)
}

type Pool struct {
	mutex *sync.Mutex
	cond  *sync.Cond

	clients     []Client
	creator     func() Client
	loaned, max int
}

func newPool(creator func() Client, max int) *Pool {
	mutex := &sync.Mutex{}
	return &Pool{
		mutex: mutex,
		cond:  sync.NewCond(mutex),

		creator: creator,
		loaned:  0,
		max:     max,
	}
}

func (p *Pool) Get() (Client, error) {
	p.mutex.Lock()
	for {
		available := len(p.clients)
		switch {
		case available < 1 && p.loaned >= p.max:
			p.cond.Wait()
		case available < 1 && p.loaned < p.max:
			p.loaned++
			p.mutex.Unlock()

			return p.creator(), nil

		case available > 0:
			var client Client
			client, p.clients = p.clients[0], p.clients[1:]
			if p.loaned < p.max {
				p.loaned++
			}
			p.mutex.Unlock()
			return client, nil
		}
	}
}

func (p *Pool) Put(client Client) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if client == nil {
		if p.loaned > 0 {
			p.loaned--
		}
		p.cond.Signal()
		return
	}

	p.clients = append(p.clients, client)

	if p.loaned > 0 {
		p.loaned--
	}
	p.cond.Signal()
}
