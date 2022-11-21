package typex

import (
	"testing"
	"testing/quick"
)

func TestWithCode(t *testing.T) {
	if BadRequest.With("Test").String() != "Bad Request Test" {
		t.Fail()
	}
}

func TestWithCode_Length(t *testing.T) {
	text := "Test with lots of characters and too many to fill in to the error code"
	if BadRequest.With(text).String() != "Bad Request Test with lots of characters and too many to fill in to the err" {
		t.Fail()
	}
}

func TestIsCode_BadRequest(t *testing.T) {
	if !BadRequest.Is(BadRequest.Name()) {
		t.Fail()
	}
}

func TestIsCode_With(t *testing.T) {
	test := "Test"
	if !BadRequest.With(test).Is(As(test)) {
		t.Fail()
	}
}

func TestIsCode_Quick(t *testing.T) {
	f := func(s string) bool {
		return BadRequest.With(s).Is(As(s))
	}
	if err := quick.Check(f, nil); err != nil {
		t.Fail()
	}
}
