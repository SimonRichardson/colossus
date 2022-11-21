package pool

import (
	"sync"

	"github.com/google/flatbuffers/go"
)

var (
	defaultPool = New(10)
)

type Pool struct {
	mutex *sync.Mutex

	items  []*flatbuffers.Builder
	loaned int
	max    int
}

func New(max int) *Pool {
	return &Pool{
		mutex: &sync.Mutex{},

		items:  []*flatbuffers.Builder{},
		loaned: 0,
		max:    max,
	}
}

func (p *Pool) Get() *flatbuffers.Builder {
	p.mutex.Lock()
	for {
		available := len(p.items)
		switch {
		case available < 1 && p.loaned >= p.max:
			// The pool was exhausted, expand!
			p.max *= 2

		case available < 1 && p.loaned < p.max:
			p.loaned++
			p.mutex.Unlock()
			return flatbuffers.NewBuilder(0)

		case available > 0:
			var builder *flatbuffers.Builder
			builder, p.items = p.items[0], p.items[1:]
			if p.loaned < p.max {
				p.loaned++
			}
			p.mutex.Unlock()
			builder.Reset()
			return builder
		}
	}
}

func (p *Pool) Put(x *flatbuffers.Builder) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if x == nil {
		if p.loaned > 0 {
			p.loaned--
		}
		return
	}

	if len(p.items) >= p.max {
		return
	}

	x.Reset()
	p.items = append(p.items, x)
	if p.loaned > 0 {
		p.loaned--
	}
}

func Get() *flatbuffers.Builder {
	return defaultPool.Get()
}

func Put(x *flatbuffers.Builder) {
	defaultPool.Put(x)
}

func SetMax(max int) {
	defaultPool.max = max
}
