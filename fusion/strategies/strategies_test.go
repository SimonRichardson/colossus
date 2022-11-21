package strategies

import (
	"flag"
	"os"
	"testing"
	"testing/quick"

	"github.com/SimonRichardson/colossus/teleprinter"
)

var (
	defaultUseStubs = false
)

func TestMain(t *testing.M) {
	var flagStubs bool
	flag.BoolVar(&flagStubs, "stubs", false, "enable stubs testing")
	flag.Parse()

	teleprinter.DefaultLog()

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

func TestHash(t *testing.T) {
	var (
		f = func(a string) int {
			return NewHash().Select(a, 10)
		}
		g = func(a string) int {
			return NewHash().Select(a, 10)
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestRoundRobin(t *testing.T) {
	var (
		f = func(a string, b int) int {
			return NewRoundRobin().Select(a, b)
		}
		g = func(a string, b int) int {
			return NewRoundRobin().Select(a, b)
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}
