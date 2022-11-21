package responses

import (
	"net/http"
	"time"

	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/schemas/pool"
	"github.com/SimonRichardson/colossus/schemas/records"
	"github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/typex"
)

func Respond(w http.ResponseWriter, status int, fn func(http.ResponseWriter), duration time.Duration) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("X-Duration", duration.String())
	w.WriteHeader(status)

	if status == http.StatusNoContent {
		return
	}

	fn(w)
}

func OKNoCotent(w http.ResponseWriter, duration time.Duration) {
	Respond(w, http.StatusNoContent, func(w http.ResponseWriter) {}, duration)
}

func OKVersion(w http.ResponseWriter, payload records.Version, duration time.Duration) {
	Respond(w, http.StatusOK, func(w http.ResponseWriter) {
		fb := pool.Get()
		defer pool.Put(fb)

		record, err := records.OKVersion{
			Duration: duration,
			Records:  payload,
		}.Write(fb)

		if err != nil {
			owned := typex.Errorf(errors.Source, typex.InternalServerError,
				"Unable to create response for OKVersion.").With(err)
			RespondError(w, "Unknown", "Unknown", typex.InternalServerError, owned)
			return
		}

		w.Write(record)
	}, duration)
}

func OKInt(w http.ResponseWriter, payload int, duration time.Duration) {
	Respond(w, http.StatusOK, func(w http.ResponseWriter) {
		fb := pool.Get()
		defer pool.Put(fb)

		record, err := records.OKInt{
			Duration: duration,
			Records:  payload,
		}.Write(fb)

		if err != nil {
			owned := typex.Errorf(errors.Source, typex.InternalServerError,
				"Unable to create response for OKInt.").With(err)
			RespondError(w, "Unknown", "Unknown", typex.InternalServerError, owned)
			return
		}

		w.Write(record)
	}, duration)
}

func OKKeyFieldScoreTxnValue(w http.ResponseWriter,
	payload selectors.KeyFieldScoreTxnValue,
	duration time.Duration,
) {
	Respond(w, http.StatusOK, func(w http.ResponseWriter) {
		fb := pool.Get()
		defer pool.Put(fb)

		record, err := records.OKKeyFieldScoreTxnValue{
			Duration: duration,
			Records:  records.FromKeyFieldScoreTxnValue(payload),
		}.Write(fb)

		if err != nil {
			owned := typex.Errorf(errors.Source, typex.InternalServerError,
				"Unable to create response for OKKeyFieldScoreTxnValue.").With(err)
			RespondError(w, "Unknown", "Unknown", typex.InternalServerError, owned)
			return
		}

		w.Write(record)
	}, duration)
}

func OKKeyFieldScoreTxnValues(w http.ResponseWriter,
	payload []selectors.KeyFieldScoreTxnValue,
	duration time.Duration,
) {
	Respond(w, http.StatusOK, func(w http.ResponseWriter) {
		fb := pool.Get()
		defer pool.Put(fb)

		record, err := records.OKKeyFieldScoreTxnValues{
			Duration: duration,
			Records:  records.FromKeyFieldScoreTxnValues(payload),
		}.Write(fb)

		if err != nil {
			owned := typex.Errorf(errors.Source, typex.InternalServerError,
				"Unable to create response for OKKeyFieldScoreTxnValues.").With(err)
			RespondError(w, "Unknown", "Unknown", typex.InternalServerError, owned)
			return
		}

		w.Write(record)
	}, duration)
}

func OKQuery(w http.ResponseWriter,
	payload []selectors.QueryRecord,
	duration time.Duration,
) {
	Respond(w, http.StatusOK, func(w http.ResponseWriter) {
		fb := pool.Get()
		defer pool.Put(fb)

		recs, err := records.FromQueryRecords(payload)
		if err != nil {
			owned := typex.Errorf(errors.Source, typex.BadRequest,
				"Unable to read QueryRecords.").With(err)
			RespondError(w, "Unknown", "Unknown", typex.BadRequest, owned)
			return
		}

		record, err := records.OKQuery{
			Duration: duration,
			Records:  recs,
		}.Write(fb)

		if err != nil {
			owned := typex.Errorf(errors.Source, typex.InternalServerError,
				"Unable to create response for OKQuery.").With(err)
			RespondError(w, "Unknown", "Unknown", typex.InternalServerError, owned)
			return
		}

		w.Write(record)
	}, duration)
}

func NoContent(w http.ResponseWriter, duration time.Duration) {
	Respond(w, http.StatusNoContent, nil, duration)
}
