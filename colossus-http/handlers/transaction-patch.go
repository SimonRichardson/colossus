package handlers

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/colossus-http/responses"
	"github.com/SimonRichardson/colossus/coordinator"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/schemas/pool"
	"github.com/SimonRichardson/colossus/schemas/schema"
	"github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/typex"
)

// TransactionPatch deletes items into the collection
func TransactionPatch(co *coordinator.Coordinator) http.HandlerFunc {
	return guard(func(w http.ResponseWriter, r *http.Request) {
		began := time.Now()

		queryKey := r.URL.Query().Get(":key")
		if !bson.IsObjectIdHex(queryKey) {
			responses.BadRequest(w, r, typex.Errorf(errors.Source, errors.InvalidArgument,
				"Invalid Key: %s", queryKey))
			return
		}

		queryId := r.URL.Query().Get(":id")
		if !bson.IsObjectIdHex(queryId) {
			responses.BadRequest(w, r, typex.Errorf(errors.Source, errors.InvalidArgument,
				"Invalid Id: %s", queryId))
			return
		}

		operations, score, maxSize, expiry, err := readPatchRecords(r.Body)
		if err != nil {
			responses.BadRequest(w, r, err)
			return
		}

		var (
			key = bs.Key(queryKey)
			id  = bs.Key(queryId)

			results, resultsErr = co.ModifyWithOperations(key,
				id,
				operations,
				score,
				selectors.SizeExpiry{
					Size:   maxSize,
					Expiry: expiry,
				})
		)
		if resultsErr != nil {
			responses.Error(w, r, resultsErr)
			return
		}

		responses.OKInt(w, results, time.Since(began))
		return
	})
}

func readPatchRecords(read io.ReadCloser) ([]selectors.Operation, float64, int64, time.Duration, error) {
	var (
		buffer bytes.Buffer
		fail   = func(err error) ([]selectors.Operation, float64, int64, time.Duration, error) {
			return nil, 0, 0, 0, err
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
		request = schema.GetRootAsPatchRequest(body, 0)

		score   = request.Score()
		maxSize = request.MaxSize()
		expiry  = request.Expiry()

		num    = request.OperationsLength()
		result = make([]selectors.Operation, num)
		fb     = pool.Get()
	)
	defer pool.Put(fb)

	for i := 0; i < num; i++ {
		operation := &schema.Operation{}
		if !request.Operations(operation, i) {
			return fail(typex.Errorf(errors.Source, errors.InvalidArgument, "Invalid Operation: %d", i))
		}

		var (
			op   = Op(string(operation.Op()))
			path = Path(string(operation.Path()))
		)

		if !op.Valid() || !path.Valid() {
			return fail(typex.Errorf(errors.Source, errors.InvalidArgument, "Invalid Operation: %d", i))
		}

		result[i] = selectors.Operation{
			Op:    selectors.Op(op.String()),
			Path:  selectors.Path(path.String()),
			Value: string(operation.Value()),
		}
	}

	return result, score, int64(maxSize), time.Duration(expiry), nil
}

type Op string

func (o Op) Valid() bool {
	switch o.String() {
	case coordinator.Match.String(), coordinator.Replace.String():
		return true
	}
	return false
}

func (o Op) String() string {
	return string(o)
}

type Path string

func (p Path) Valid() bool {
	if str := p.String(); len(str) > 1 && str[0:1] == "/" {
		return true
	}
	return false
}

func (p Path) String() string {
	return string(p)
}
