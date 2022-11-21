package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"testing/quick"
	"time"

	"github.com/SimonRichardson/colossus/colossus-shim/common"
	"github.com/SimonRichardson/colossus/colossus-shim/coordinator"
	"github.com/SimonRichardson/colossus/colossus-shim/env"
	"github.com/SimonRichardson/colossus/quatsch"
	b "github.com/SimonRichardson/colossus/quatsch/pool/bson"
	"github.com/SimonRichardson/colossus/tests"
	"github.com/SimonRichardson/colossus/typex"
)

var (
	defaultUseStubs = false
)

func TestMain(t *testing.M) {
	var flagStubs bool
	flag.BoolVar(&flagStubs, "stubs", false, "enable stubs testing")
	flag.Parse()

	defaultUseStubs = flagStubs

	res := t.Run()
	tearDown()
	os.Exit(res)
}

func setup(e *env.Env) (*httptest.Server, *coordinator.Coordinator) {
	e.C.Instrumentation = "Noop"
	e.C.Logs = "Noop"

	server := newServer(e)
	return httptest.NewServer(server.Handler), server.co
}

func tear(ts *httptest.Server) {
	ts.Close()
}

func tearDown() {

}

func getIdentPool() quatsch.Pool {
	var (
		maxBuffer               = 99999
		maxInsertionPerDuration = int64(1000000)
	)

	return quatsch.New(b.New(maxBuffer, time.Second, maxInsertionPerDuration))
}

// Test Version

type VersionResponse struct {
	Duration int64  `json:"duration"`
	Records  string `json:"records"`
}

func testVersion(url string,
	co *coordinator.Coordinator,
) (func() string, func() string) {
	var (
		f = func() string {
			body, err := common.Get(fmt.Sprintf("%s/http/version", url), func(headers http.Header) {
				headers.Set("Accept", "application/json")
				headers.Set("Content-Type", "application/json")
			})

			if err != nil {
				typex.Fatal(err)
			}

			var resp VersionResponse
			tests.MustUnmarshal(body, &resp)

			return resp.Records
		}
		g = func() string {
			return "0.0.1"
		}
	)
	return f, g
}

func TestVersion(t *testing.T) {
	e := env.New(nil)

	ts, co := setup(e)
	defer tear(ts)

	f, g := testVersion(ts.URL, co)

	if err := quick.CheckEqual(f, g, tests.Config()); err != nil {
		t.Error(err)
	}
}

// Test Reserve

func testReserve(url string,
	co *coordinator.Coordinator,
) func() bool {
	pool := getIdentPool()
	return func() bool {
		key, err1 := b.Bson(pool.Get())
		if err1 != nil {
			typex.Fatal(err1)
		}

		var (
			uri       = fmt.Sprintf("%s/events/%s/tickets/reserve/%d", url, key.Hex(), 2)
			body, err = common.Post(uri, []byte{}, func(headers http.Header) {
				headers.Set("Accept", "application/json")
				headers.Set("Content-Type", "application/json")
			})
		)
		if err != nil {
			typex.Fatal(err)
		}
		return len(body) > 1
	}
}

func TestReserve(t *testing.T) {
	e := env.New(nil)
	e.C.HttpAddress = "http://127.0.0.1:9002"

	ts, co := setup(e)
	defer tear(ts)

	f := testReserve(ts.URL, co)

	if err := quick.Check(f, tests.Config()); err != nil {
		t.Error(err)
	}
}
