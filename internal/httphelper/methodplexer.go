package httphelper

import "net/http"

type MethodPlexer struct {
	Get  http.HandlerFunc
	Post http.HandlerFunc
	Put  http.HandlerFunc
	// Feel free to add more
}

func MethodPlexMiddleware(p MethodPlexer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			p.Get(w, r)
		case http.MethodPost:
			p.Post(w, r)
		case http.MethodPut:
			p.Put(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
