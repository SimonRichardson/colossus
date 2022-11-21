package strategies

import (
	"sync"

	"github.com/SimonRichardson/colossus/fusion"
	"github.com/SimonRichardson/colossus/fusion/selectors"
	"github.com/SimonRichardson/colossus/teleprinter"
)

type keyStore struct {
	mutex    sync.Mutex
	store    selectors.KeyStore
	prefix   selectors.Prefix
	fallback fusion.SelectionStrategy
	values   map[string]int
	ticker   chan struct{}
}

// NewKeyStore defines a way to get select a node from a remote server, useful
// for determining which node a key should alway go to!
func NewKeyStore(prefix selectors.Prefix,
	ticker chan struct{},
	store selectors.KeyStore,
) *keyStore {
	k := &keyStore{
		mutex:    sync.Mutex{},
		prefix:   prefix,
		store:    store,
		fallback: NewHash(),
		values:   make(map[string]int),
		ticker:   ticker,
	}
	go k.run()
	return k
}

func (r *keyStore) Select(key string, max int) int {
	r.mutex.Lock()

	if v, ok := r.values[key]; ok {
		r.mutex.Unlock()
		return v
	}
	r.mutex.Unlock()

	return r.fallback.Select(key, max)
}

func (r *keyStore) run() {
	for range r.ticker {
		if values, err := r.store.List(r.prefix); err != nil {
			teleprinter.L.Error().Printf("Unable to get store list : %s", err.Error())
		} else {
			r.mutex.Lock()
			r.values = values
			r.mutex.Unlock()
		}
	}
}
