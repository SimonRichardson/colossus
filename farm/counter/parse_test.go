package counter

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestParseFarmString(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	for farmString, expected := range map[string]struct {
		success     bool
		numClusters int
	}{
		"":                                                {false, 0}, // no entries
		";;;":                                             {false, 0}, // no entries
		"foo1:1234":                                       {true, 1},
		"foo1:1234;bar1:1234":                             {true, 2},
		"foo1:1234;;bar1:1234":                            {false, 0}, // empty middle cluster
		"foo1,writeonly":                                  {false, 0}, // writeonly is an invalid token now
		"a1:1234,a2:1234;b1:1234,b2:1234":                 {true, 2},
		"a1:1234,a2:1234; b1:1234,b2:1234 ":               {true, 2},
		"a1:1234,a2:1234; b1:1234,b2:1234; ":              {false, 0}, // empty last cluster
		"a1:1234,a2:1234;b1:1234,b2:1234,writeonly":       {false, 0}, // writeonly is an invalid token now
		"a1:1234,a2:1234,a3:1234;b1:1234,b2:1234,b3:1234": {true, 2},
		"a1:1234,a2:1234 ; b1:1234,b2:1234 ; c1:1234":     {true, 3},
	} {
		clusters, err := ParseString(
			farmString,
			"1s", "1s", "1s",
			"RoundRobin",
			1,
			nil,
		)
		if expected.success && err != nil {
			t.Errorf("%q: %s", farmString, err)
			continue
		}
		if !expected.success && err == nil {
			t.Errorf("%q: expected error, got none", farmString)
			continue
		}
		if expected, got := expected.numClusters, len(clusters); expected != got {
			t.Errorf("expected %d cluster(s), got %d", expected, got)
		}
	}
}
