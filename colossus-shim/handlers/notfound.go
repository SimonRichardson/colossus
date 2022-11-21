package handlers

import (
	"net/http"

	"github.com/SimonRichardson/colossus/colossus-shim/responses"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/typex"
)

func NotFound() func(http.ResponseWriter, *http.Request) {
	return handle(func(w http.ResponseWriter, r *http.Request) {
		responses.NotFound(w, r, typex.Errorf(errors.Source, errors.MissingContent,
			"Not Found"))
		return
	})
}
