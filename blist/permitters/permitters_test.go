package permitters

import (
	"flag"
	"os"
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

func TestPermitter_AllowedWithMore(t *testing.T) {
	var (
		f = func(s int64) bool {
			n := abs(s)
			p := New(n, time.Second)
			return p.Allowed(n + 1)
		}
		g = func(s int64) bool {
			return false
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestPermitter_AllowedWithLess(t *testing.T) {
	var (
		f = func(s int64) bool {
			n := abs(s)
			p := New(n, time.Second)
			return p.Allowed(n - 1)
		}
		g = func(s int64) bool {
			return true
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestPermitter_AllowAny(t *testing.T) {
	var (
		f = func(a, b int64) bool {
			n := -abs(a)
			p := New(n, time.Second)
			return p.Allowed(b)
		}
		g = func(a, b int64) bool {
			return true
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

// Throttle

func TestPermitter_ThrottleAllowedWithMore(t *testing.T) {
	var (
		f = func(s int64) bool {
			n := abs(s)
			p := NewThrottle(n, time.Millisecond)
			return p.Allowed(n + 1)
		}
		g = func(s int64) bool {
			return true
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func abs(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}
