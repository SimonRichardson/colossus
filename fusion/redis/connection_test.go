package redis

import (
	"io/ioutil"
	"log"
	"math"
	"runtime"
	"testing"
	"time"

	"github.com/SimonRichardson/colossus/teleprinter"
)

func TestMemoryRegression(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	teleprinter.DefaultLog()

	var (
		quit = make(chan struct{})
		diff = make(chan uint64)
		tick = time.Tick(250 * time.Millisecond)
	)

	go func() {
		var (
			m        runtime.MemStats
			biggest  uint64
			smallest = uint64(math.MaxUint64)
		)
		for {
			select {
			case <-tick:
				runtime.ReadMemStats(&m)
				if m.HeapAlloc > biggest {
					biggest = m.HeapAlloc
				}
				if m.HeapAlloc < smallest {
					smallest = m.HeapAlloc
				}
			case <-quit:
				diff <- biggest - smallest
				return
			}
		}
	}()

	var (
		addr           = "127.0.0.1:54321" // invalid
		timeout        = 500 * time.Millisecond
		maxConnections = 25
	)
	p := newConnectionPool(addr, "", newConnectionTimeout().All(timeout), maxConnections, nil)
	for i, n := 0, 10; i < n; i++ {
		runtime.GC()
		p.get()
	}

	close(quit)
	if delta := <-diff; delta > 100 {
		t.Errorf("HeapAlloc ∆ was %d", delta)
	}
}
