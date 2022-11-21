package mongo

import (
	"testing"
	"testing/quick"
	"time"
)

func TestTimeout_All(t *testing.T) {
	var (
		f = func(d time.Duration) time.Duration {
			return newConnectionTimeout().All(d).global
		}
		g = func(d time.Duration) time.Duration {
			return d
		}
	)
	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}
