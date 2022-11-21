package handlers

import (
	"net/http"

	"github.com/SimonRichardson/colossus/teleprinter"
)

func handle(fn func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		teleprinter.L.Info().Printf("Requesting %s.\n", r.URL)
		fn(w, r)
	}
}
