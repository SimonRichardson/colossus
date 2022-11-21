package alertmanager

import (
	"sync"
	"time"
)

const (
	defaultDurationInterval = time.Second
	defaultDurationOffset   = defaultDurationInterval * 2
)

type node struct {
	id  uint64
	msg string
}

type cancellable struct {
	base AlertBase
	id   uint64
}

func (c cancellable) Cancel() {
	for _, nodes := range c.base.nodes {
	inner:
		for k, v := range nodes {
			if v.id == c.id {
				nodes = append(nodes[:k], nodes[k+1:]...)
				break inner
			}
		}
	}
}

type AlertBase struct {
	mutex   *sync.Mutex
	counter uint64
	nodes   map[time.Time][]node
}

func Make(fn func(string)) AlertBase {
	b := AlertBase{
		mutex:   &sync.Mutex{},
		counter: 0,
		nodes:   make(map[time.Time][]node),
	}
	go b.run(fn)
	return b
}

func (m AlertBase) Dispatch(msg string, options AlertOptions) Cancellable {
	m.counter++

	var (
		now   = time.Now()
		id    = m.counter
		delay = options.Delay
	)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i := uint(0); i < options.Repititions; i++ {
		key := clamp(now.Add((delay * time.Duration(i)) + defaultDurationOffset))
		m.nodes[key] = append(m.nodes[key], node{
			id:  id,
			msg: msg,
		})
	}

	return cancellable{
		base: m,
		id:   id,
	}
}

func (m AlertBase) run(fn func(string)) {
	for t := range time.Tick(defaultDurationInterval) {
		m.mutex.Lock()

		if nodes, ok := m.nodes[t]; ok {
			for _, v := range nodes {
				fn(v.msg)
			}
		}

		// clean up any nodes left
		for k := range m.nodes {
			if k.Before(t) {
				delete(m.nodes, k)
			}
		}

		m.mutex.Unlock()
	}
}

func clamp(t time.Time) time.Time {
	// Traps the time to 1 second
	return time.Unix(t.Unix(), 0)
}
