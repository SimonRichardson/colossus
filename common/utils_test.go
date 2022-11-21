package common

import (
	"errors"
	"testing"
)

func TestNormalise(t *testing.T) {
	inOuts := [][]string{
		[]string{"basic", "basic"},
		[]string{"BaSic", "basic"},
		[]string{"BaS ic", "basic"},
		[]string{"BaS ic ", "basic"},
		[]string{`BaS i
      c`, "basic"},
		[]string{`
      `, ""},
		[]string{` `, ""},
	}

	for _, v := range inOuts {
		if Normalise(v[0]) != v[1] {
			t.Errorf("%v(%v) != %v", "Normalise", v[0], v[1])
			t.FailNow()
		}
	}
}

func TestStripWhitespace(t *testing.T) {
	inOuts := [][]string{
		[]string{"basic", "basic"},
		[]string{"BaSic", "BaSic"},
		[]string{"BaS ic", "BaSic"},
		[]string{"BaS ic ", "BaSic"},
		[]string{`BaS i
      c`, "BaSic"},
		[]string{`
      `, ""},
		[]string{` `, ""},
	}

	for _, v := range inOuts {
		if StripWhitespace(v[0]) != v[1] {
			t.Errorf("%v(%v) != %v", "StripWhitespace", v[0], v[1])
			t.FailNow()
		}
	}
}

// "testing/quick" pointless here as would be just rewriting function
type ReduceIntsIO struct {
	Ins []int
	Out int
}

func TestSumIntegers(t *testing.T) {
	inOuts := []ReduceIntsIO{
		ReduceIntsIO{
			Ins: []int{1, 2, 3},
			Out: 6,
		},
		ReduceIntsIO{
			Ins: []int{-1, 2, 3},
			Out: 4,
		},
		ReduceIntsIO{
			Ins: []int{},
			Out: 0,
		},
		ReduceIntsIO{
			Ins: []int{3},
			Out: 3,
		},
		ReduceIntsIO{
			Ins: []int{-1, -2, -3},
			Out: -6,
		},
	}

	for _, v := range inOuts {
		if SumIntegers(v.Ins) != v.Out {
			t.Errorf("%v(%v) != %v", "SumIntegers", v.Ins, v.Out)
			t.FailNow()
		}
	}
}

func TestAvgIntegers(t *testing.T) {
	inOuts := []ReduceIntsIO{
		ReduceIntsIO{
			Ins: []int{1, 2, 3},
			Out: 2,
		},
		ReduceIntsIO{
			Ins: []int{-1, 2, 3},
			Out: 1,
		},
		ReduceIntsIO{
			Ins: []int{},
			Out: 0,
		},
		ReduceIntsIO{
			Ins: []int{3},
			Out: 3,
		},
		ReduceIntsIO{
			Ins: []int{-1, -2, -3},
			Out: -2,
		},
	}

	for _, v := range inOuts {
		if AvgIntegers(v.Ins) != v.Out {
			t.Errorf("%v(%v) != %v", "AvgIntegers", v.Ins, v.Out)
			t.FailNow()
		}
	}
}

func TestFilterErrors(t *testing.T) {
	ins := []error{
		errors.New("My Error"),
		nil,
		errors.New("My second error")}
	outs := FilterErrors(ins)

	if len(outs) != 2 {
		t.Errorf("len(%v(%v)) != 2", "FilterErrors", ins)
		t.FailNow()
	}

	for _, v := range outs {
		if v == nil {
			t.Errorf("%v(%v) contained a nil error: %v", "FilterErrors", ins, outs)
			t.FailNow()
		}
	}
}

type ReduceErrorsIO struct {
	Ins []error
	Out error
}

func TestSumErrors(t *testing.T) {
	ins := []ReduceErrorsIO{
		ReduceErrorsIO{
			Ins: []error{errors.New("Hello "), errors.New("World!")},
			Out: errors.New("Hello ; World!"),
		},
		ReduceErrorsIO{
			Ins: []error{errors.New("Hello "), nil, errors.New("World!")},
			Out: errors.New("Hello ; World!"),
		},
		ReduceErrorsIO{
			Ins: []error{},
			Out: errors.New(""),
		},
	}

	for _, v := range ins {
		if SumErrors(v.Ins).Error() != v.Out.Error() {
			t.Errorf("%v(%v) != %v", "SumErrors", v.Ins, v.Out)
			t.FailNow()
		}
	}
}
