package common

import (
	"time"
)

func Tick(d time.Duration) chan struct{} {
	ch := make(chan struct{})

	go func() {
		tick := time.Tick(d)
		for range tick {
			ch <- struct{}{}
		}
	}()

	return ch
}
