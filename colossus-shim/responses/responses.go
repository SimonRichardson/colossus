package responses

import (
	"encoding/json"
	"net/http"
	"time"
)

func ResponseWithEmptyBody(w http.ResponseWriter, status int, duration time.Duration) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("x-duration", duration.String())
	w.WriteHeader(status)
}

func Respond(w http.ResponseWriter, status int, records interface{}, duration time.Duration) {
	ResponseWithEmptyBody(w, status, duration)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"records":  records,
		"duration": duration,
	})
}

func OKWithBytes(w http.ResponseWriter, bytes []byte, duration time.Duration) {
	ResponseWithEmptyBody(w, http.StatusOK, duration)

	w.Write(bytes)
}

func OK(w http.ResponseWriter, records interface{}, duration time.Duration) {
	Respond(w, http.StatusOK, records, duration)
}

func OKWithNoRecords(w http.ResponseWriter, duration time.Duration) {
	Respond(w, http.StatusOK, map[string]interface{}{}, duration)
}

func NoContent(w http.ResponseWriter, duration time.Duration) {
	ResponseWithEmptyBody(w, http.StatusNoContent, duration)
}
