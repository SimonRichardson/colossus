package handlers

import (
	"fmt"
	"net/http"

	"time"

	"strconv"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	c "github.com/SimonRichardson/colossus/colossus-shim/common"
	"github.com/SimonRichardson/colossus/colossus-shim/coordinator"
	"github.com/SimonRichardson/colossus/colossus-shim/responses"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/schemas/records"
	"github.com/SimonRichardson/colossus/typex"
	flatbuffers "github.com/google/flatbuffers/go"
	"gopkg.in/mgo.v2/bson"
)

var (
	defaultTime = time.Date(2016, 12, 1, 1, 1, 1, 1, time.UTC)
)

func Reserve(co *coordinator.Coordinator, host string) http.HandlerFunc {
	return handle(func(w http.ResponseWriter, r *http.Request) {
		var (
			err   error
			score int

			began  = time.Now()
			amount = 0

			query       = r.URL.Query()
			queryKey    = query.Get(":key")
			queryAmount = query.Get(":amount")
		)
		if !bson.IsObjectIdHex(queryKey) {
			responses.BadRequest(w, r, typex.Errorf(errors.Source, errors.InvalidArgument,
				"Invalid Key: %s", queryKey))
			return
		}
		if amount, err = strconv.Atoi(queryAmount); err != nil {
			responses.BadRequest(w, r, typex.Errorf(errors.Source, errors.InvalidArgument,
				"Invalid Amount: %s", queryAmount))
			return
		}

		score, err = co.Increment(bs.Key(queryKey), defaultTime)
		if err != nil {
			responses.InternalServerError(w, r, typex.Errorf(errors.Source, errors.Fatal,
				"Error: %s", err.Error()))
			return
		}

		var (
			values = makePostRecords(amount)
			record = records.PostRecords{
				Records: values,
				Score:   float64(score),
				MaxSize: 99999, // TODO - find this out, do we care!?
				Expiry:  time.Minute * 5,
			}
			bytes, writeErr = record.Write(flatbuffers.NewBuilder(0))
		)
		if writeErr != nil {
			responses.InternalServerError(w, r, typex.Errorf(errors.Source, errors.Fatal,
				"Error: %s", err.Error()))
			return
		}

		bytes, err = c.Post(fmt.Sprintf("%s/http/v1/%s", host, queryKey), bytes, func(headers http.Header) {
			headers.Set("Accept", "application/octet-stream")
			headers.Set("Content-Type", "application/octet-stream")
		})

		if err != nil {
			responses.BadRequest(w, r, err)
			return
		}

		responses.OKWithBytes(w, bytes, time.Since(began))
		return
	})
}

func makePostRecords(amount int) []records.PostRecord {
	var (
		now   = time.Now()
		owner = bson.NewObjectId()
		txn   = bson.NewObjectId()

		res = make([]records.PostRecord, 0, amount)
	)

	for i := 0; i < amount; i++ {
		res = append(res, records.PostRecord{
			Id:       bson.NewObjectId(),
			Updated:  now,
			Reserved: now,
			Expiry:   now.Add(time.Minute * 5),
			Cost: records.Cost{
				Currency: "GBP",
				Price:    uint64(0),
			},
			OwnerId:       owner,
			TransactionId: txn,
		})
	}

	return res
}
