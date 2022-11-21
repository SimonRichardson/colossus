package responses

import (
	"net/http"

	"fmt"

	"github.com/SimonRichardson/colossus/schemas/pool"
	"github.com/SimonRichardson/colossus/schemas/records"
	"github.com/SimonRichardson/colossus/teleprinter"
	"github.com/SimonRichardson/colossus/typex"
)

const defaultShowInspectTrace bool = false

func RespondError(w http.ResponseWriter,
	method, url string,
	fallback typex.ErrorCode,
	err error,
) {
	var (
		code    = getCode(err, fallback)
		errHttp = fmt.Errorf("%s %s: HTTP %d", method, url, code)
		errFull = typex.Lift(err).With(errHttp)
		trace   = typex.Inspect(errFull)
	)

	teleprinter.L.Error().Println(trace)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(code)

	record := records.Error{
		Error:       err.Error(),
		Code:        getName(err, code),
		Description: getDescription(trace),
	}

	fb := pool.Get()
	defer pool.Put(fb)

	bytes, err := record.Write(fb)
	if err != nil {
		panic(err)
	}

	w.Write(bytes)
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	InternalServerError(w, r, err)
}

func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	RespondError(w, r.Method, r.URL.String(), typex.BadRequest, err)
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	RespondError(w, r.Method, r.URL.String(), typex.InternalServerError, err)
}

func NotFound(w http.ResponseWriter, r *http.Request, err error) {
	RespondError(w, r.Method, r.URL.String(), typex.NotFound, err)
}

func getCode(err error, fallback typex.ErrorCode) int {
	if code := typex.ErrCode(err); code >= 0 {
		return code
	}
	return fallback.HTTPStatusCode()
}

func getName(err error, code int) string {
	if name := typex.ErrName(err); name != "" {
		return name
	}
	return http.StatusText(code)
}

func getDescription(trace string) string {
	if defaultShowInspectTrace {
		return trace
	}
	return ""
}
