package httphelper

import (
	"encoding/json"
	"net/http"

	"golang.org/x/exp/slog"
)

func StructResponse[T any](handler func(r *http.Request) (T,*HTTPError) ) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
	response, hErr := handler(r)
	if hErr != nil{
		slog.With(
			"status",hErr.Code,
			"cause", hErr.Cause.Error(),
			"message",hErr.Message	,
		).Error("http error")
		w.WriteHeader(hErr.Code)
		_,err := w.Write([]byte(hErr.Message))
		if err != nil {
			slog.With(
				"error", err.Error(),
				"message",hErr.Message,).
				Error("unable to write error message")
		}
		return
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil{
		slog.With(			"error",err.Error()).
		Error("unable to write json response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	}
}