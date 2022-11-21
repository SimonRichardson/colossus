package handlers

import "net/http"

func Charge() http.HandlerFunc {
	return handle(func(w http.ResponseWriter, r *http.Request) {
	})
}
