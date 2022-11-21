package cache

import "testing"
import "time"

func TestService(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fail()
		}
	}()
	DefaultService("tcp://lorenz", time.Minute)
}
