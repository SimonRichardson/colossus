package tests

import (
	"math/rand"
	"testing"
	"testing/quick"
	"time"
)

func Config() *quick.Config {
	if testing.Short() {
		return &quick.Config{
			MaxCount:      10,
			MaxCountScale: 10,
			Rand:          rand.New(rand.NewSource(time.Now().UnixNano())),
		}
	}
	return &quick.Config{
		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func ConfigFor(amount int) *quick.Config {
	return &quick.Config{
		MaxCount:      amount,
		MaxCountScale: float64(amount),
		Rand:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}
