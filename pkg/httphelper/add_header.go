package httphelper

import (
	"net/http"
)

func AddHeader(key, value string, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		w.Header().Add(key, value)
	}
}
