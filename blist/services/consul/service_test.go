package consul

import "testing"

func TestService(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fail()
		}
	}()
	DefaultService("http://localhost/consul", "blist", "output")
}
