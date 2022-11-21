package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/SimonRichardson/colossus/colossus-http/handlers"
	"github.com/SimonRichardson/colossus/coordinator"
	"github.com/SimonRichardson/colossus/env"
	"github.com/SimonRichardson/colossus/errors"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/teleprinter"
	"github.com/SimonRichardson/colossus/teleprinter/logs/parse"
	"github.com/SimonRichardson/colossus/tests"
	"github.com/SimonRichardson/colossus/typex"
	"github.com/gorilla/pat"
	"gopkg.in/mgo.v2/bson"
)

var (
	defaultRandom = rand.New(rand.NewSource(time.Now().UnixNano()))
	defaultExpiry = time.Minute * 5
)

func main() {
	var (
		iterationCountPtr = flag.Int("itercount", 4, "How many iterations to do in total")
		scaleFactorPtr    = flag.Int("scalefactor", 10, "Scale factor of the total number of requests")
		totalRequestsPtr  = flag.Int("totalreqs", 100, "Total number of requests to send")
		ticketCountPtr    = flag.Int("ticketcount", 2, "Number of tickets to allocate each request")
		recordPtr         = flag.Bool("record", true, "If recording should be enabled or not")
		ciServerPtr       = flag.Bool("ciserver", false, "If the ci server is running the app")
		clusterServerPtr  = flag.Bool("clusterserver", false, "If the cluster server is running the app")
		verbosePtr        = flag.Bool("verbose", true, "Verbose mode")
	)

	flag.Parse()

	var (
		iterationCount = *iterationCountPtr
		scaleFactor    = *scaleFactorPtr
		totalRequests  = *totalRequestsPtr
		ticketCount    = *ticketCountPtr
		record         = *recordPtr
		ciServer       = *ciServerPtr
		clusterServer  = *clusterServerPtr
		verbose        = *verbosePtr

		csv = newCSV()
	)

	log.SetFlags(0)

	if !verbose {
		log.SetOutput(ioutil.Discard)
	}

	for i := 0; i < iterationCount; i++ {
		var (
			total = totalRequests * scale(scaleFactor, i)
			began = time.Now()
		)

		log.Printf("=== RUN   PerformanceRequest_%d_%d\n", total, ticketCount)

		duration, result := performanceRequest(total, ticketCount)

		csv.add(total, ticketCount, result, duration, ciServer, clusterServer)

		log.Printf("--- PASS: PerformanceRequest_%d_%d (%s)\n", total, ticketCount, time.Since(began).String())
	}

	result := csv.write()

	if record {
		upload(result)
	}

	log.Println("-----------------------\n| Performance Results |\n-----------------------")
	log.Printf(string(result))
	log.Println("-----------------------")
}

func scale(fact, exp int) (val int) {
	val = 1
	for i := 0; i < exp; i++ {
		val *= fact
	}
	return
}

func performanceRequest(total, amount int) (time.Duration, int) {
	e := env.New(nil)
	e.Logs = "Noop"
	e.Instrumentation = "Noop"

	var (
		ts, co = setup(e)

		key = bson.NewObjectId()

		inserter = api(ts.URL, key, (total*amount)*2, co)
		bodies   = make([]tests.PostBody, total)
	)

	for i := 0; i < total; i++ {
		bodies[i] = generatePostBody(amount)
	}

	defer tear(ts)

	began := time.Now()

	// Do this so we don't get optimized out!
	var response int
	for i := 0; i < total; i++ {
		response = inserter(bodies[i])
	}

	return time.Since(began), response
}

func generatePostBody(amount int) tests.PostBody {
	return tests.PostBody(nil).Make(defaultRandom, amount)
}

func setupLogging(e *env.Env) {
	var err error
	if teleprinter.L, err = parse.ParseString(e.Logs); err != nil {
		typex.Fatal(err)
	}
}

type server struct {
	HttpAddress string
	Handler     http.Handler
	co          *coordinator.Coordinator
}

func newServer(e *env.Env) server {
	// Setup logging
	setupLogging(e)

	var (
		co = coordinator.New(e, func(value s.KeyFieldScoreTxnValue) (map[string]interface{}, error) {
			return nil, typex.Errorf(errors.Source, errors.UnexpectedResults, "No transformer")
		}, accessor{})

		path = func(p string) func(string) string {
			return func(n string) string { return fmt.Sprintf("%s%s", p, n) }
		}

		prefix  = path("/http/v1")
		tprefix = path(prefix("/{key}"))
		router  = pat.New()
	)

	router.Post(tprefix(""), handlers.TransactionsPost(co))

	return server{
		e.HttpAddress,
		http.Handler(router),
		co,
	}
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

type accessor struct{}

func (a accessor) GetFieldValue(i interface{}, field string) (string, error) {
	return "", nil
}

func (a accessor) SetFieldValue(i interface{}, field, value string) error {
	return nil
}
