package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/SimonRichardson/colossus/colossus-http/handlers"
	"github.com/SimonRichardson/colossus/common"
	"github.com/SimonRichardson/colossus/coordinator"
	"github.com/SimonRichardson/colossus/env"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/schemas/pool"
	"github.com/SimonRichardson/colossus/schemas/records"
	"github.com/SimonRichardson/colossus/schemas/schema"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/teleprinter"
	"github.com/SimonRichardson/colossus/teleprinter/logs/parse"
	"github.com/SimonRichardson/colossus/typex"
	"github.com/gorilla/pat"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultTimeAfterReload = time.Second
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
	pool.SetMax(1000)

	var (
		e      = env.New(nil)
		server = newServer(e)

		accessor = coordinator.NewCoordinatorAccessor(server.co)
		alerts   = accessor.AlertManager()
	)

	// hot reloading
	go func() {
		watcher := e.Watch()
		for {
			select {
			case <-watcher:
				// Pause the server
				func() {
					server.co.Pause()
					defer server.co.Resume()

					if err := server.co.Topology(e); err != nil {
						typex.Fatal("PANIC: Disaster!", err)

						alerts.TopologyPanic()
					}

					time.Sleep(defaultTimeAfterReload)
				}()
			}
		}
	}()

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

	router.Get("/http/version", handlers.Version(e.Version))

	router.Get(tprefix("/query"), handlers.TransactionsQuery(co))
	router.Get(tprefix("/count"), handlers.TransactionsCount(co))
	router.Delete(tprefix("/rollback"), handlers.TransactionsRollback(co))

	// Transaction
	// The following are handlers for doing individual requests and
	// transformations on a set (collection) of transactions

	router.Get(tprefix("/{id}"), handlers.TransactionGet(co))
	router.Patch(tprefix("/{id}"), handlers.TransactionPatch(co))

	// Transactions
	// The following are handlers for doing bulk requests and transformations
	// on a set (collection) of transactions

	router.Get(tprefix(""), handlers.TransactionsGet(co))
	router.Post(tprefix(""), handlers.TransactionsPost(co))
	router.Put(tprefix(""), handlers.TransactionsPut(co))
	router.Delete(tprefix(""), handlers.TransactionsDelete(co))

	// Custom verbs.
	router.Add("COUNT", tprefix(""), handlers.TransactionsCount(co))
	router.Add("QUERY", tprefix(""), handlers.TransactionsQuery(co))
	router.Add("ROLLBACK", tprefix(""), handlers.TransactionsRollback(co))

	router.NotFoundHandler = http.HandlerFunc(handlers.NotFound())

	return server{
		e.HttpAddress,
		http.Handler(router),
		co,
	}
}

type accessor struct{}

func (a accessor) GetFieldValue(i interface{}, field string) (string, error) {
	record, ok := i.(*records.PostRecord)
	if !ok {
		return "", typex.Errorf(errors.Source, errors.InvalidArgument, "Invalid Type")
	}

	switch field {
	case "txn":
		return record.TransactionId.Hex(), nil
	case "owner_id":
		return record.OwnerId.Hex(), nil
	case "expiry_time":
		return fmt.Sprintf("%d", record.Expiry.UnixNano()), nil
	default:
		return "", typex.Errorf(errors.Source, errors.UnexpectedResults,
			"Invalid property %s", field)
	}
}

func (a accessor) SetFieldValue(i interface{}, field, value string) error {
	record, ok := i.(*records.PostRecord)
	if !ok {
		return typex.Errorf(errors.Source, errors.InvalidArgument, "Invalid Type")
	}

	switch field {
	case "txn":
		if !bson.IsObjectIdHex(value) {
			return typex.Errorf(errors.Source, errors.InvalidArgument,
				"Invalid ObjectId (%s)", value)
		}
		record.TransactionId = bson.ObjectIdHex(value)
	case "owner_id":
		if !bson.IsObjectIdHex(value) {
			return typex.Errorf(errors.Source, errors.InvalidArgument,
				"Invalid ObjectId (%s)", value)
		}
		record.OwnerId = bson.ObjectIdHex(value)
	case "expiry_time":
		s, err := strconv.Atoi(value)
		if err != nil {
			return typex.Errorf(errors.Source, errors.InvalidArgument,
				"Invalid time (%s)", value)
		}
		record.Expiry = time.Unix(0, int64(s))
	default:
		return typex.Errorf(errors.Source, errors.InvalidArgument,
			"Invalid property %s", field)
	}
	return nil
}

func transformer(value s.KeyFieldScoreTxnValue) (map[string]interface{}, error) {
	header, err := records.ReadType(value.Value)
	if err != nil {
		return nil, err
	}

	var (
		meta = func(updated time.Time) map[string]interface{} {
			m := map[string]interface{}{
				"model": map[string]interface{}{
					"created_at": time.Now(),
					"updated_at": updated,
				},
			}
			return m
		}
		m = map[string]interface{}{}
	)

	switch header {
	case schema.TypePost:
		var (
			record    = &records.PostRecord{}
			body, err = records.ReadBody(value.Value)
		)
		if err != nil {
			return nil, err
		}
		if err = record.Read(body); err != nil {
			return nil, err
		}

		m["_id"] = bson.ObjectIdHex(value.Field.String())
		m["owner_id"] = record.OwnerId
		m["expiry_time"] = record.Expiry
		m["reserved_at"] = record.Reserved
		m["meta"] = meta(record.Updated)
		m["txn"] = record.TransactionId

		cost := record.Cost
		m["cost"] = map[string]interface{}{
			"currency": cost.Currency,
			"price":    cost.Price,
		}

	case schema.TypePut:
		var (
			record    = &records.PutRecord{}
			body, err = records.ReadBody(value.Value)
		)
		if err != nil {
			return nil, err
		}
		if err = record.Read(body); err != nil {
			return nil, err
		}

		m["_id"] = bson.ObjectIdHex(value.Field.String())
		m["owner_id"] = record.OwnerId
		m["purchased_at"] = record.Purchased
		m["meta"] = meta(record.Updated)
		m["txn"] = record.TransactionId

		cost := record.EventCost
		m["cost"] = map[string]interface{}{
			"currency": cost.Currency,
			"price":    cost.Price,
		}

		dates := record.EventDates
		m["event_date"] = time.Unix(0, int64(dates.Start))
		m["event_date_end"] = time.Unix(0, int64(dates.End))

		codes := record.Codes
		m["bar_code"] = map[string]interface{}{
			"type":   codes.BarcodeType,
			"origin": codes.BarcodeOrigin,
			"source": codes.BarcodeSource,
		}
		m["qr_code"] = codes.QRCode

	default:
		return nil, typex.Errorf(errors.Source, errors.NoCaseFound, "Unknown Type")
	}

	return m, nil
}
