package bson

import (
	"sync"
	"time"

	"github.com/SimonRichardson/colossus/quatsch/pool"
	"gopkg.in/mgo.v2/bson"
)

const (
	defaultTimeout           = time.Second
	defaultIterationDuration = defaultTimeout / 3
)

type Pool struct {
	mutex       *sync.Mutex
	max, amount int
	dst         chan bson.ObjectId
	duration    time.Duration
	tactic      pool.Tactic
}

func New(max int, duration time.Duration, perDuration int64) *Pool {
	pool := &Pool{
		mutex:    &sync.Mutex{},
		max:      max,
		dst:      make(chan bson.ObjectId, max),
		duration: duration,
		tactic:   pool.GetTactic(duration, perDuration),
	}

	// First generate a healthy buffer
	pool.generate()
	// Keep generating new ones all the time.
	go pool.iterate()

	return pool
}

func (p *Pool) get() (bson.ObjectId, error) {
	timer := time.NewTimer(defaultTimeout)
	for {
		select {
		case id := <-p.dst:
			p.amount--
			timer.Stop()
			return id, nil

		case <-timer.C:
			return bson.ObjectId(""), pool.ErrExhaustion
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

func (p *Pool) WithTyped(fn func(bson.ObjectId) error) error {
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
			// Make sure that we never get any duplications.
			res := make(map[bson.ObjectId]struct{}, diff)
			for i := 0; i < diff; i++ {
				id := bson.NewObjectId()
				if _, ok := res[id]; ok {
					// TODO: Should log out about a collision occuring!
					continue
				}
				res[id] = struct{}{}
			}

			for k := range res {
				p.amount++
				p.dst <- k
			}
			return nil
		}); err != nil {
			diff /= 2
			continue
		}

		break
	}
}

func Bson(r interface{}, err error) (bson.ObjectId, error) {
	if err != nil {
		return bson.ObjectId(""), err
	}

	if x, ok := r.(bson.ObjectId); ok {
		return x, nil
	}

	return bson.ObjectId(""), pool.ErrTypeError
}
