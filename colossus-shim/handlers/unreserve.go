package handlers

import "net/http"

func Unreserve() http.HandlerFunc {
	return handle(func(w http.ResponseWriter, r *http.Request) {
	})
}
