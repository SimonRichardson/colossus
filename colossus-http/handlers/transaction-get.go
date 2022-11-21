package handlers

import (
	"net/http"
	"time"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/colossus-http/responses"
	"github.com/SimonRichardson/colossus/coordinator"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/typex"
	"gopkg.in/mgo.v2/bson"
)

func TransactionGet(co *coordinator.Coordinator) http.HandlerFunc {
	return accepts(func(w http.ResponseWriter, r *http.Request) {
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

		var (
			key = bs.Key(queryKey)
			id  = bs.Key(queryId)

			results, resultsErr = co.Select(key, id)
		)
		if resultsErr != nil {
			responses.Error(w, r, resultsErr)
			return
		}

		responses.OKKeyFieldScoreTxnValue(w, results, time.Since(began))
		return
	})
}
