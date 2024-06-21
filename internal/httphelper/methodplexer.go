package httphelper

import (
	"log/slog"
	"net/http"
)

type MethodPlexer struct {
	Get  http.HandlerFunc
	Post http.HandlerFunc
	Put  http.HandlerFunc
	// Feel free to add more
}

func MethodPlexMiddleware(p MethodPlexer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet:
			p.Get(w, r)
		case r.Method == http.MethodPost && p.Post != nil:
			p.Post(w, r)
		case r.Method == http.MethodPut && p.Put != nil:
			p.Put(w, r)
		default:
			slog.With(
				"method", r.Method,
				"path", r.URL.Path,
			).DebugContext(r.Context(), "Method not allowed")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
