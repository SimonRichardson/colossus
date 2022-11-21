package handlers

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/colossus-http/responses"
	"github.com/SimonRichardson/colossus/coordinator"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/typex"
)

// TransactionsGet represents the end point for selecting items from the store.
func TransactionsGet(co *coordinator.Coordinator) http.HandlerFunc {
	return accepts(func(w http.ResponseWriter, r *http.Request) {
		began := time.Now()

		queryKey := r.URL.Query().Get(":key")
		if !bson.IsObjectIdHex(queryKey) {
			responses.BadRequest(w, r, typex.Errorf(errors.Source, errors.InvalidArgument,
				"Invalid Key: %s", queryKey))
			return
		}

		var (
			maxSize, ok0    = parseInt(r.Form, "size", 10)
			expiry, ok1     = parseInt(r.Form, "expiry", 10)
			queryLimit, ok2 = parseInt(r.Form, "limit", 10)
		)

		if !ok0 || !ok1 || !ok2 {
			responses.BadRequest(w, r, typex.Errorf(errors.Source, errors.InvalidArgument,
				"Invalid request parameter"))
			return
		}

		var (
			key          = bs.Key(queryKey)
			sizeExpiry   = selectors.MakeKeySizeSingleton(key, int64(maxSize), time.Duration(expiry))
			results, err = co.SelectRange(key, queryLimit, sizeExpiry)
		)
		if err != nil {
			responses.InternalServerError(w, r, err)
			return
		}

		responses.OKKeyFieldScoreTxnValues(w, results, time.Since(began))
		return
	})
}

func parseInt(values url.Values, key string, defaultValue int) (int, bool) {
	valueStr := values.Get(key)
	if valueStr == "" {
		return defaultValue, false
	}
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return defaultValue, false
	}
	return int(value), true
}
