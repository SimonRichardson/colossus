package mongo

import (
	"testing"
	"testing/quick"
	"time"
)

func TestParse_RandomStrategy(t *testing.T) {
	var (
		f = func(d time.Duration) bool {
			_, _, err := Parse(d.String(), "random")
			return err == nil
		}
	)
	if err := quick.Check(f, config()); err != nil {
		t.Error(err)
	}
}

func TestParse_HashStrategy(t *testing.T) {
	var (
		f = func(d time.Duration) bool {
			_, _, err := Parse(d.String(), "hash")
			return err == nil
		}
	)
	if err := quick.Check(f, config()); err != nil {
		t.Error(err)
	}
}

func TestParse_RoundRobinStrategy(t *testing.T) {
	var (
		f = func(d time.Duration) bool {
			_, _, err := Parse(d.String(), "roundrobin")
			return err == nil
		}
	)
	if err := quick.Check(f, config()); err != nil {
		t.Error(err)
	}
}

func TestParse_NoiseStrategy(t *testing.T) {
	var (
		f = func(d time.Duration, s string) bool {
			_, _, err := Parse(d.String(), s)
			return err != nil
		}
	)
	if err := quick.Check(f, config()); err != nil {
		t.Error(err)
	}
}
