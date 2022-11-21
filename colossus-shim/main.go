package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/SimonRichardson/colossus/colossus-shim/coordinator"
	"github.com/SimonRichardson/colossus/colossus-shim/env"
	"github.com/SimonRichardson/colossus/colossus-shim/handlers"
	"github.com/SimonRichardson/colossus/common"
	"github.com/SimonRichardson/colossus/teleprinter"
	"github.com/SimonRichardson/colossus/teleprinter/logs/parse"
	"github.com/SimonRichardson/colossus/typex"
	"github.com/gorilla/pat"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type server struct {
	HttpAddress string
	Handler     http.Handler
	co          *coordinator.Coordinator
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	rand.Seed(time.Now().UnixNano())

	var (
		e      = env.New(nil)
		server = newServer(e)
	)

	log.Printf("listening on %s", server.HttpAddress)
	typex.Fatal(common.ListenAndServe(
		server.HttpAddress,
		common.ServerTimeout{
			Read:  e.C.HttpReadTimeout,
			Write: e.C.HttpWriteTimeout,
		},
		teleprinter.L.Error(),
		server.Handler,
		server.co.Quit,
	))
}

func setupLogging(e *env.Env) {
	var err error
	if teleprinter.L, err = parse.ParseString(e.C.Logs); err != nil {
		typex.Fatal(err)
	}
}

func newServer(e *env.Env) server {
	// Setup logging
	setupLogging(e)

	var (
		co = coordinator.New(e)

		router = pat.New()
		host   = e.C.HttpAddress
	)

	// Order of these are fundamental!

	if e.C.PrometheusMetrics {
		router.Handle("/metrics", promhttp.Handler())
	}

	router.Get("/http/version", handlers.Version(e.C.Version))

	router.Post("/events/{key}/tickets/reserve/{amount}", handlers.Reserve(co, host))
	router.Post("/events/{key}/tickets/unreserve", handlers.Unreserve())
	router.Post("/events/{key}/tickets/charge", handlers.Charge())

	router.NotFoundHandler = http.HandlerFunc(handlers.NotFound())

	return server{
		e.ShimHttpAddress,
		http.Handler(router),
		co,
	}
}
