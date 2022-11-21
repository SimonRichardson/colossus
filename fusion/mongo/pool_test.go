package mongo

import (
	"testing"
	"testing/quick"

	"github.com/SimonRichardson/colossus/fusion/strategies"
	mgo "gopkg.in/mgo.v2"
)

func TestPool_Size(t *testing.T) {
	var (
		f = func(a []string, b int) int {
			pool := New(a, strategies.NewRandom(), newConnectionTimeout(), b, func(*mgo.DialInfo) (Session, error) {
				return &session{}, nil
			})
			return pool.Size()
		}
		g = func(a []string, b int) int {
			return len(a)
		}
	)
	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestPool_Index(t *testing.T) {
	var (
		pool = New([]string{}, strategies.NewRandom(), newConnectionTimeout(), 100, func(*mgo.DialInfo) (Session, error) {
			return &session{}, nil
		})
		f = func(a string) int {
			return pool.Index(a)
		}
		g = func(a string) int {
			return pool.Index(a)
		}
	)
	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestPool_With(t *testing.T) {
	var (
		f = func(a []string, b string) (res bool) {
			// Skip the test!
			if len(a) < 1 {
				res = true
				return
			}

			pool := New(a, strategies.NewRandom(), newConnectionTimeout(), 100, func(*mgo.DialInfo) (Session, error) {
				return &session{}, nil
			})
			if err := pool.With(b, func(s Session) error {
				res = true
				return nil
			}); err != nil {
				t.Fatal(err)
			}
			return
		}
	)
	if err := quick.Check(f, config()); err != nil {
		t.Error(err)
	}
}

func TestPool_WithIndex(t *testing.T) {
	var (
		f = func(a []string, b int) (res bool) {
			// Skip the test!
			if len(a) < 1 || b < 1 {
				res = true
				return
			}

			pool := New(a, strategies.NewRandom(), newConnectionTimeout(), 100, func(*mgo.DialInfo) (Session, error) {
				return &session{}, nil
			})
			if err := pool.WithIndex(b%len(a), func(s Session) error {
				res = true
				return nil
			}); err != nil {
				t.Fatal(err)
			}
			return
		}
	)
	if err := quick.Check(f, config()); err != nil {
		t.Error(err)
	}
}
