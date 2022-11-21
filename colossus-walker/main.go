package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/SimonRichardson/colossus/colossus-http/handlers"
	"github.com/SimonRichardson/colossus/colossus-walker/agents"
	"github.com/SimonRichardson/colossus/common"
	"github.com/SimonRichardson/colossus/coordinator"
	"github.com/SimonRichardson/colossus/env"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/schemas/pool"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/teleprinter"
	"github.com/SimonRichardson/colossus/teleprinter/logs/parse"
	"github.com/SimonRichardson/colossus/typex"
	"github.com/gorilla/pat"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultVersion = "0.0.1"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lmicroseconds)

	rand.Seed(time.Now().UnixNano())
	pool.SetMax(1000)

	var (
		e      = env.New(nil)
		server = newServer(e)
	)

	if err := server.Daemon(); err != nil {
		typex.Fatalf("Error starting Daemon supervisor, with : %s\n", err.Error())
	}

	log.Printf("listening on %s", server.HttpAddress)
	typex.Fatal(common.ListenAndServe(
		server.HttpAddress,
		common.ServerTimeout{
			Read:  e.HttpReadTimeout,
			Write: e.HttpWriteTimeout,
		},
		teleprinter.L.Error(),
		server.Handler,
		server.co.Quit,
	))
}

type server struct {
	HttpAddress string
	Handler     http.Handler
	co          *coordinator.Coordinator
	agents      []agents.Agent
}

func (s server) Daemon() error {
	// TODO: We should supervise the agents!
	opts := agents.AgentOptions{
		Coordinator: s.co,
		HttpAddress: s.HttpAddress,
	}

	for _, v := range s.agents {
		if err := v.Init(opts); err != nil {
			return err
		}
	}
	return nil
}

func setupLogging(e *env.Env) {
	var err error
	if teleprinter.L, err = parse.ParseString(e.Logs); err != nil {
		typex.Fatal(err)
	}
}

func newServer(e *env.Env) server {
	// Setup logging
	setupLogging(e)

	var (
		co = coordinator.New(e, transformer, accessor{})

		path = func(p string) func(string) string {
			return func(n string) string { return fmt.Sprintf("%s%s", p, n) }
		}

		prefix  = path("/http/v1")
		tprefix = path(prefix("/{key}"))
		router  = pat.New()
	)

	// Order of these are fundamental!

	if e.PrometheusMetrics {
		router.Handle("/metrics", promhttp.Handler())
	}

	router.Get("/http/version", handlers.Version(defaultVersion))
	router.Get(tprefix("/select"), handlers.TransactionsGet(co))

	router.NotFoundHandler = http.HandlerFunc(handlers.NotFound())

	return server{
		e.HttpAddress,
		http.Handler(router),
		co,
		[]agents.Agent{
			agents.Walk{},
		},
	}
}

type accessor struct{}

func (a accessor) GetFieldValue(interface{}, string) (string, error) {
	return "", fmt.Errorf("missing implementation.")
}
func (a accessor) SetFieldValue(interface{}, string, string) error {
	return fmt.Errorf("missing implementation.")
}

func transformer(s.KeyFieldScoreTxnValue) (map[string]interface{}, error) {
	return nil, typex.Errorf(errors.Source, errors.UnexpectedResults, "No transformer")
}
