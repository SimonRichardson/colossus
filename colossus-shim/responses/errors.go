package responses

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SimonRichardson/colossus/teleprinter"
	"github.com/SimonRichardson/colossus/typex"
)

const defaultShowInspectTrace bool = false

func RespondError(w http.ResponseWriter,
	method, url string,
	fallback typex.ErrorCode,
	err error,
	fn func(error, int, string) (map[string]interface{}, bool),
) {
	var (
		code    = getCode(err, fallback)
		errHttp = fmt.Errorf("%s %s: HTTP %d", method, url, code)
		errFull = typex.Lift(err).With(errHttp)
		trace   = typex.Inspect(errFull)
	)

	teleprinter.L.Error().Println(trace)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if res, ok := fn(errFull, code, trace); ok {
		json.NewEncoder(w).Encode(res)
	}
}

func RespondErrorWithEmptyBody(w http.ResponseWriter,
	method, url string,
	fallback typex.ErrorCode,
	err error,
) {
	RespondError(w, method, url, fallback, err, func(error, int, string) (map[string]interface{}, bool) {
		return nil, false
	})
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	InternalServerError(w, r, err)
}

func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	RespondError(w, r.Method, r.URL.String(), typex.BadRequest, err, generic)
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	RespondError(w, r.Method, r.URL.String(), typex.InternalServerError, err, generic)
}

func NotFound(w http.ResponseWriter, r *http.Request, err error) {
	RespondError(w, r.Method, r.URL.String(), typex.NotFound, err, generic)
}

func generic(err error, code int, description string) (map[string]interface{}, bool) {
	return map[string]interface{}{
		"error":       err.Error(),
		"code":        getName(err, code),
		"description": getDescription(description),
	}, true
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
