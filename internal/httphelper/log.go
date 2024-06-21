package httphelper

import (
	"net/http"

	"log/slog"
)

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.With(
			"method", r.Method,
			"path", r.URL.Path,
		).DebugContext(r.Context(), "request")
		handler.ServeHTTP(w, r)
	})
}
