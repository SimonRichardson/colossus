package mongo

import (
	"flag"
	"math"
	"os"
	"runtime"
	"testing"
	"testing/quick"
	"time"

	"github.com/SimonRichardson/colossus/teleprinter"
	mgo "gopkg.in/mgo.v2"
)

var (
	defaultUseStubs     = false
	defaultMongoAddress = ""
)

func TestMain(t *testing.M) {
	var flagStubs bool
	flag.BoolVar(&flagStubs, "stubs", false, "enable stubs testing")
	flag.Parse()

	teleprinter.DefaultLog()

	defaultUseStubs = flagStubs

	// FIXME: Locate this from the env vars
	defaultMongoAddress = "mongo:27017"

	os.Exit(t.Run())
}

func config() *quick.Config {
	if testing.Short() {
		return &quick.Config{
			MaxCount:      10,
			MaxCountScale: 10,
		}
	}
	return nil
}

type session struct{}

func (s *session) Ping() error             { return nil }
func (s *session) Copy() Session           { return s }
func (s *session) Close()                  {}
func (s *session) DB(name string) Database { return nil }

func TestConnectionPool_Get(t *testing.T) {
	var (
		f = func(a string) bool {
			_, err := NewConnectionPool(a, newConnectionTimeout(), 100, func(*mgo.DialInfo) (Session, error) {
				return &session{}, nil
			}).Get()
			return err == nil
		}
	)
	if err := quick.Check(f, config()); err != nil {
		t.Error(err)
	}
}

func TestConnectionPool_GetMultipleTimes(t *testing.T) {
	var (
		conn = NewConnectionPool("add", newConnectionTimeout(), 10000, func(*mgo.DialInfo) (Session, error) {
			return &session{}, nil
		})
		f = func(a string) bool {
			_, err := conn.Get()
			return err == nil
		}
	)
	if err := quick.Check(f, config()); err != nil {
		t.Error(err)
	}
}

func TestConnectionPool_GetPut(t *testing.T) {
	var (
		f = func(a string) bool {
			conn := NewConnectionPool(a, newConnectionTimeout(), 1, func(*mgo.DialInfo) (Session, error) {
				return &session{}, nil
			})
			s, err := conn.Get()
			if err != nil {
				t.Fatal(err)
			}
			conn.Put(s)
			s, err = conn.Get()
			return err == nil
		}
	)
	if err := quick.Check(f, config()); err != nil {
		t.Error(err)
	}
}

func TestMemoryRegression(t *testing.T) {
	if defaultUseStubs || testing.Short() {
		t.Skip()
	}

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
		timeout        = 500 * time.Millisecond
		maxConnections = 25
	)
	p := NewConnectionPool(defaultMongoAddress, newConnectionTimeout().All(timeout), maxConnections, nil)
	for i, n := 0, 10; i < n; i++ {
		runtime.GC()
		p.Get()
	}

	close(quit)
	if delta := <-diff; delta > 100 {
		t.Errorf("HeapAlloc ∆ was %d", delta)
	}
}
