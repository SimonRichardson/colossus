package handlers

import (
	"bytes"
	"io"
	"net/http"
	"time"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/colossus-http/responses"
	"github.com/SimonRichardson/colossus/coordinator"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/schemas/pool"
	"github.com/SimonRichardson/colossus/schemas/records"
	"github.com/SimonRichardson/colossus/schemas/schema"
	"github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/typex"
	"gopkg.in/mgo.v2/bson"
)

// TransactionsDelete deletes items into the collection
func TransactionsDelete(co *coordinator.Coordinator) http.HandlerFunc {
	return guard(func(w http.ResponseWriter, r *http.Request) {
		began := time.Now()

		queryKey := r.URL.Query().Get(":key")
		if !bson.IsObjectIdHex(queryKey) {
			responses.BadRequest(w, r, typex.Errorf(errors.Source, errors.InvalidArgument,
				"Invalid Key: %s", queryKey))
			return
		}

		fieldValues, score, maxSize, expiry, err := readDeleteRecords(r.Body)
		if err != nil {
			responses.BadRequest(w, r, err)
			return
		}

		var (
			key           = bs.Key(queryKey)
			maxSizeExpiry = selectors.MakeKeySizeSingleton(key, maxSize, expiry)

			elements           = fieldValues.KeyFieldScoreTxnValues(key, score)
			results, deleteErr = co.Delete(elements, maxSizeExpiry)
		)
		if deleteErr != nil {
			responses.Error(w, r, deleteErr)
			return
		}

		responses.OKInt(w, results, time.Since(began))
		return
	})
}

func readDeleteRecords(read io.ReadCloser) (selectors.FieldTxnValues, float64, int64, time.Duration, error) {
	var (
		buffer bytes.Buffer
		fail   = func(err error) (selectors.FieldTxnValues, float64, int64, time.Duration, error) {
			return nil, 0, 0, time.Duration(0), err
		}
	)
	if _, err := buffer.ReadFrom(read); err != nil {
		return fail(typex.Errorf(errors.Source, errors.InvalidArgument,
			"Invalid Body"))
	}

	body := buffer.Bytes()
	if len(body) < 1 {
		return fail(typex.Errorf(errors.Source, errors.InvalidArgument,
			"Invalid Body Length"))
	}

	var (
		request = schema.GetRootAsDeleteRequest(body, 0)
		score   = request.Score()
		maxSize = request.MaxSize()
		expiry  = request.Expiry()
	)
	if maxSize < 1 {
		return fail(typex.Errorf(errors.Source, errors.InvalidArgument, "Invalid MaxSize: %d", maxSize))
	}
	if expiry < 1 {
		return fail(typex.Errorf(errors.Source, errors.InvalidArgument, "Invalid expiry: %d", expiry))
	}

	var (
		num    = request.RecordsLength()
		result = make([]selectors.FieldTxnValue, num)
		fb     = pool.Get()
	)
	defer pool.Put(fb)

	for i := 0; i < num; i++ {
		record := &schema.DeleteRecord{}
		if !request.Records(record, i) {
			return fail(typex.Errorf(errors.Source, errors.InvalidArgument, "Invalid Record: %d", i))
		}

		id, err := readRecordId(record)
		if err != nil {
			return fail(err)
		}

		transaction, err := readRecordTransactionId(record)
		if err != nil {
			return fail(err)
		}

		fb.Reset()

		value, err := records.DeleteRecordFromSchemaToByte(fb, record)
		if err != nil {
			return fail(err)
		}
		result[i] = selectors.FieldTxnValue{
			Field: id,
			Txn:   transaction,
			Value: records.PackageDeleteRecord(value),
		}
	}

	return result, score, int64(maxSize), time.Duration(expiry), nil
}
