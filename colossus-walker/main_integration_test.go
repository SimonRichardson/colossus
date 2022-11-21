package main

import (
	"flag"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"testing/quick"

	"github.com/SimonRichardson/colossus/colossus-walker/common"
	"github.com/SimonRichardson/colossus/coordinator"
	"github.com/SimonRichardson/colossus/env"
	"github.com/SimonRichardson/colossus/schemas/records"
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

	os.Exit(t.Run())
}

func setup(e *env.Env) (*httptest.Server, *coordinator.Coordinator) {
	e.Logs = "Noop"
	e.Instrumentation = "Noop"

	server := newServer(e)
	return httptest.NewServer(server.Handler), server.co
}

func tear(ts *httptest.Server) {
	ts.Close()
}

func config() *quick.Config {
	if testing.Short() {
		return &quick.Config{
			MaxCount:      10,
			MaxCountScale: 10,
		}
	}
	return nil
}

// Test Version

func testVersion(url string,
	co *coordinator.Coordinator,
) (func() string, func() string) {
	var (
		f = func() string {
			body, err := common.Get(fmt.Sprintf("%s/http/version", url))
			if err != nil {
				typex.Fatal(err)
			}

			s := &records.OKVersion{}
			s.Read(body)

			return s.Records.Version
		}
		g = func() string {
			return defaultVersion
		}
	)
	return f, g
}

func TestVersion(t *testing.T) {
	e := env.New(nil)

	ts, co := setup(e)
	defer tear(ts)

	f, g := testVersion(ts.URL, co)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}
