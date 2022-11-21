package tests

import (
	"math/rand"
	"sync"
	"time"
)

var (
	nowNano = time.Now().UnixNano()
	Random  = rand.New(&lockedSource{src: rand.NewSource(nowNano)})
)

type lockedSource struct {
	mutex sync.Mutex
	src   rand.Source
}

func (r *lockedSource) Int63() int64 {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.src.Int63()
}

func (r *lockedSource) Seed(seed int64) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.src.Seed(seed)
}
