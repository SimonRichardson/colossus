package handlers

import (
	"net/http"
	"time"

	"github.com/SimonRichardson/colossus/colossus-http/responses"
	"github.com/SimonRichardson/colossus/common"
	"github.com/SimonRichardson/colossus/schemas/records"
)

func Version(version string) http.HandlerFunc {
	parseErr := common.ParseSemver(version)

	return handle(func(w http.ResponseWriter, r *http.Request) {
		began := time.Now()

		if parseErr != nil {
			responses.InternalServerError(w, r, parseErr)
			return
		}

		responses.OKVersion(w, records.Version{version}, time.Since(began))
		return
	})
}
