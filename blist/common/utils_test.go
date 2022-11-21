package common

import "testing"

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
