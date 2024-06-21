package httphelper

import (
	"encoding/json"
	"net/http"

	"log/slog"
)

func StructResponse[T any](handler func(r *http.Request) (T, *HTTPError)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// I think it makes sense to give the entire request to the handler. In general handlers are
		// varied in what they need from the request, be it path params, query params, body, header, ...
		// and it doesn't make sense to abstract that away because in the end it would just recreate http.Request.
		response, hErr := handler(r)
		if hErr != nil {
			slog.With(
				"status", hErr.Code,
				"cause", hErr.Cause.Error(),
				"message", hErr.Message,
			).ErrorContext(r.Context(), "http error")
			w.WriteHeader(hErr.Code)
			_, err := w.Write([]byte(hErr.Message))
			if err != nil {
				slog.With(
					"error", err.Error(),
					"message", hErr.Message).
					ErrorContext(r.Context(), "unable to write error message")
			}
			return
		}
		w.Header().Add("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			slog.With("error", err.Error()).
				ErrorContext(r.Context(), "unable to write json response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
