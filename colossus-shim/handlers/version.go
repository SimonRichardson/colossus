package handlers

import (
	"net/http"
	"time"

	"github.com/SimonRichardson/colossus/colossus-shim/responses"
	"github.com/SimonRichardson/colossus/common"
)

func Version(version string) http.HandlerFunc {
	parseErr := common.ParseSemver(version)

	return handle(func(w http.ResponseWriter, r *http.Request) {
		began := time.Now()

		if parseErr != nil {
			responses.InternalServerError(w, r, parseErr)
			return
		}

		responses.OK(w, version, time.Since(began))
		return
	})
}
