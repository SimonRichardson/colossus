package handlers

import (
	"net/http"
	"time"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/colossus-http/responses"
	"github.com/SimonRichardson/colossus/coordinator"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/typex"
	"gopkg.in/mgo.v2/bson"
)

// TransactionsQuery queries items in the collection
func TransactionsQuery(co *coordinator.Coordinator) http.HandlerFunc {
	return accepts(func(w http.ResponseWriter, r *http.Request) {
		began := time.Now()

		queryKey := r.URL.Query().Get(":key")
		if !bson.IsObjectIdHex(queryKey) {
			responses.BadRequest(w, r, typex.Errorf(errors.Source, errors.InvalidArgument,
				"Invalid Key: %s", queryKey))
			return
		}

		var (
			maxSize, ok0 = parseInt(r.Form, "size", 10)
			expiry, ok1  = parseInt(r.Form, "expiry", 10)
		)

		if !ok0 || !ok1 {
			responses.BadRequest(w, r, typex.Errorf(errors.Source, errors.InvalidArgument, "Invalid request parameter"))
			return
		}

		ownerId := r.Form.Get("owner_id")
		if !bson.IsObjectIdHex(ownerId) {
			responses.BadRequest(w, r, typex.Errorf(errors.Source, errors.InvalidArgument, "Invalid request parameter"))
			return
		}

		var (
			key                 = bs.Key(queryKey)
			results, resultsErr = co.Query(key, selectors.QueryOptions{
				OwnerId: bs.Key(ownerId),
			}, selectors.SizeExpiry{
				Size:   int64(maxSize),
				Expiry: time.Duration(expiry),
			})
		)
		if resultsErr != nil {
			responses.InternalServerError(w, r, resultsErr)
			return
		}

		responses.OKQuery(w, results, time.Since(began))
		return
	})
}
