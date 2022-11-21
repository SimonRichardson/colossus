package uuid

import (
	"sync"
	"time"

	"github.com/SimonRichardson/colossus/quatsch/pool"
)

const (
	defaultTimeout           = time.Second
	defaultIterationDuration = defaultTimeout / 3
)

type Pool struct {
	mutex       *sync.Mutex
	max, amount int
	dst         chan UUID
	duration    time.Duration
	tactic      pool.Tactic
	generator   *generator
}

func New(gen *generator, max int, duration time.Duration, perDuration int64) *Pool {
	pool := &Pool{
		mutex:     &sync.Mutex{},
		max:       max,
		dst:       make(chan UUID, max),
		duration:  duration,
		tactic:    pool.GetTactic(duration, perDuration),
		generator: gen,
	}

	// First generate a healthy buffer
	pool.generate()
	// Keep generating new ones all the time.
	go pool.iterate()

	return pool
}

func (p *Pool) get() (UUID, error) {
	timer := time.NewTimer(defaultTimeout)
	for {
		select {
		case id := <-p.dst:
			p.amount--
			timer.Stop()
			return id, nil

		case <-timer.C:
			return UUID([24]byte{}), pool.ErrExhaustion
		}
	}
}

func (p *Pool) With(fn func(interface{}) error) error {
	v, err := p.get()
	if err != nil {
		return err
	}

	return fn(v)
}

func (p *Pool) Len() int {
	return p.amount
}

func (p *Pool) iterate() {
	ticker := time.NewTicker(defaultIterationDuration)
	for {
		select {
		case <-ticker.C:
			p.generate()
		}
	}
}

func (p *Pool) generate() {
	diff := p.max - p.amount
	for i := 0; i < 3; i++ {
		if err := p.tactic(diff, func() error {
			for i := 0; i < diff; i++ {
				p.amount++
				p.dst <- p.generator.Next()
			}
			return nil
		}); err != nil {
			diff /= 2
			continue
		}

		break
	}
}
