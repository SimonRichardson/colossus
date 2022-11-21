package semaphore

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"testing"
	"testing/quick"
	"time"
)

var (
	defaultUseStubs = false
)

func TestMain(t *testing.M) {
	var flagStubs bool
	flag.BoolVar(&flagStubs, "stubs", false, "enable stubs testing")
	flag.Parse()

	defaultUseStubs = flagStubs

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

type stringBulkItem struct {
	value string
}

func (s stringBulkItem) Key() string { return s.value }
func (s stringBulkItem) Len() int    { return 1 }
func (s stringBulkItem) Equals(b BulkItem) bool {
	return s.Key() == b.Key()
}

type stringsBulkItem struct {
	values []string
}

func (s stringsBulkItem) Key() string { return strings.Join(s.values, "") }
func (s stringsBulkItem) Len() int    { return len(s.values) }
func (s stringsBulkItem) Equals(b BulkItem) bool {
	return s.Key() == b.Key()
}

// Queue

func TestQueueAdd_MultipleTimes(t *testing.T) {
	var (
		f = func(s []string) int {
			q := newQueue()
			for _, v := range s {
				if err := q.Add(stringBulkItem{v}); err != nil {
					t.Fatal(err)
				}
			}
			return q.Len()
		}
		g = func(s []string) int {
			return len(s)
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestQueueAdd_MultipleItems(t *testing.T) {
	var (
		f = func(s []string) int {
			q := newQueue()
			for range s {
				if err := q.Add(stringsBulkItem{s}); err != nil {
					t.Fatal(err)
				}
			}
			return q.Len()
		}
		g = func(s []string) int {
			return len(s) * len(s)
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestQueueRemove(t *testing.T) {
	var (
		f = func(s []string) int {
			q := newQueue()
			for _, v := range s {
				if err := q.Add(stringBulkItem{v}); err != nil {
					t.Fatal(err)
				}
			}
			if len(s) > 0 {
				index := rand.Intn(len(s))
				if err := q.Remove(stringBulkItem{s[index]}); err != nil {
					t.Fatal(err)
				}
			}
			return q.Len()
		}
		g = func(s []string) int {
			if len(s) < 1 {
				return 0
			}
			return len(s) - 1
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestQueuePeek(t *testing.T) {
	var (
		f = func(s []string) (int, int) {
			q := newQueue()
			for _, v := range s {
				if err := q.Add(stringBulkItem{v}); err != nil {
					t.Fatal(err)
				}
			}
			return len(q.Peek(len(s) / 2)), q.Len()
		}
		g = func(s []string) (int, int) {
			a := len(s) / 2
			return a, len(s) - a
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestQueuePeek_TillExhaustion(t *testing.T) {
	var (
		f = func(s []string) int {
			q := newQueue()
			for _, v := range s {
				if err := q.Add(stringBulkItem{v}); err != nil {
					t.Fatal(err)
				}
			}
			for range s {
				q.Peek(1)
			}
			return q.Len()
		}
		g = func(s []string) int {
			return 0
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestQueuePeek_TillOverExhaustion(t *testing.T) {
	var (
		f = func(s []string) int {
			q := newQueue()
			for _, v := range s {
				if err := q.Add(stringBulkItem{v}); err != nil {
					t.Fatal(err)
				}
			}
			for range s {
				q.Peek(1)
			}
			q.Peek(1)
			return q.Len()
		}
		g = func(s []string) int {
			return 0
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

// Bulk

func TestBulkAdd_WhenHitsSizeWillSendMessage(t *testing.T) {
	var (
		f = func(s []string) (res int) {
			q := NewBulk(func(items []BulkItem) {
				res = len(items)
			}, len(s)-1, time.Minute)
			for _, v := range s {
				q.Add(stringBulkItem{v})
			}
			runtime.Gosched()
			time.Sleep(time.Millisecond)
			return
		}
		g = func(s []string) int {
			return len(s)
		}
	)
	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestBulkAdd_WhenDoesNotHitSizeWillNotSendMessage(t *testing.T) {
	var (
		f = func(s []string) (res int) {
			q := NewBulk(func(items []BulkItem) {
				res = len(items)
			}, len(s), time.Minute)
			for _, v := range s {
				q.Add(stringBulkItem{v})
			}
			runtime.Gosched()
			time.Sleep(time.Millisecond)
			return
		}
		g = func(s []string) int {
			return 0
		}
	)
	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestBulkRemove_WhenHitsSizeWillSendMessage(t *testing.T) {
	var (
		f = func(s []string) (res int) {
			q := NewBulk(func(items []BulkItem) {
				res = len(items)
			}, len(s), time.Minute)
			for _, v := range s {
				q.Add(stringBulkItem{v})
			}
			if len(s) < 1 {
				res = 0
				return
			}
			q.Remove(stringBulkItem{s[0]})
			for i := 0; i < 2; i++ {
				q.Add(stringBulkItem{fmt.Sprintf("%d", i)})
			}
			runtime.Gosched()
			time.Sleep(time.Millisecond)
			return
		}
		g = func(s []string) int {
			if len(s) < 1 {
				return 0
			}
			return len(s) + 1
		}
	)
	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestBulkPeek(t *testing.T) {
	var (
		f = func(s []string) int {
			q := NewBulk(func(items []BulkItem) {
				t.Fatal("Unexpected call.")
			}, len(s), time.Minute)
			for _, v := range s {
				q.Add(stringBulkItem{v})
			}
			return len(q.Peek())
		}
		g = func(s []string) int {
			return len(s)
		}
	)
	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestBulkPeek_ThenAddSendsMessage(t *testing.T) {
	var (
		f = func(a []string) (res int) {
			q := NewBulk(func(items []BulkItem) {
				res = len(items)
			}, len(a), time.Minute)

			for _, v := range a {
				q.Add(stringBulkItem{v})
			}
			q.Peek()
			for _, v := range a {
				q.Add(stringBulkItem{v})
			}
			q.Add(stringBulkItem{fmt.Sprintf("pad-%d", rand.Int())})
			runtime.Gosched()
			time.Sleep(time.Millisecond)
			return
		}
		g = func(a []string) int {
			return len(a) + 1
		}
	)
	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

// Async time tests

func TestBulkAdd_WhenHitsTimeoutWillSendMessage(t *testing.T) {
	var (
		f = func(s []string) (res int) {
			q := NewBulk(func(items []BulkItem) {
				res = len(items)
			}, len(s), time.Millisecond)
			for _, v := range s {
				q.Add(stringBulkItem{v})
			}
			runtime.Gosched()
			time.Sleep(time.Millisecond * 10)
			q.Stop()
			return q.queue.Len()
		}
		g = func(s []string) int {
			return 0
		}
	)
	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}
