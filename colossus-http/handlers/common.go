package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/colossus-http/responses"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/schemas/schema"
	"github.com/SimonRichardson/colossus/teleprinter"
	"github.com/SimonRichardson/colossus/typex"
	"gopkg.in/mgo.v2/bson"
)

const (
	contentType = "application/octet-stream"
)

func handle(fn func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		teleprinter.L.Info().Printf("Requesting %s.\n", r.URL)
		fn(w, r)
	}
}

func accepts(fn func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return handle(func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Header.Get("Accept")) != contentType {
			responses.BadRequest(w, r, fmt.Errorf("Invalid Content-Type"))
			return
		}

		if err := r.ParseForm(); err != nil {
			responses.BadRequest(w, r, err)
			return
		}

		fn(w, r)
	})
}

func guard(fn func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return handle(func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Header.Get("Accept")) != contentType {
			responses.BadRequest(w, r, fmt.Errorf("Invalid Accept"))
			return
		}

		if strings.ToLower(r.Header.Get("Content-Type")) != contentType {
			responses.BadRequest(w, r, fmt.Errorf("Invalid Content-Type"))
			return
		}

		defer func() {
			io.Copy(ioutil.Discard, r.Body)
			r.Body.Close()
		}()

		if err := r.ParseForm(); err != nil {
			responses.BadRequest(w, r, err)
			return
		}

		fn(w, r)
	})
}

type recordWithId interface {
	Id(obj *schema.Id) *schema.Id
}

type recordWithTransactionId interface {
	TransactionId(obj *schema.Id) *schema.Id
}

func readRecordId(record recordWithId) (bs.Key, error) {
	var hex string
	if id := record.Id(nil); id != nil {
		hex = string(id.Hex())
	} else {
		return bs.Key(""), typex.Errorf(errors.Source, errors.UnexpectedResults,
			"Invalid Record Id: %v", id)
	}

	if !bson.IsObjectIdHex(hex) {
		return bs.Key(""), typex.Errorf(errors.Source, errors.UnexpectedResults,
			"Invalid Record Id (expected bson.ObjectId): %s", hex)
	}
	return bs.Key(hex), nil
}

type wrapTransactionId struct {
	record recordWithTransactionId
}

func (w wrapTransactionId) Id(obj *schema.Id) *schema.Id {
	return w.record.TransactionId(obj)
}

func readRecordTransactionId(record recordWithTransactionId) (bs.Key, error) {
	return readRecordId(wrapTransactionId{record})
}
