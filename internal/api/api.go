package api

import (
	"database/sql"
	"net/http"

	"github.com/vankleefjim/go_experiment_jet/internal/config"
	"github.com/vankleefjim/go_experiment_jet/internal/db"
	"github.com/vankleefjim/go_experiment_jet/internal/httphelper"
	"github.com/vankleefjim/go_experiment_jet/internal/todos"

	"log/slog"
)

func Routes(cfg config.Server, conn *sql.DB) *http.ServeMux {
	// This extra level is not needed when writing a microservice
	mux := http.NewServeMux()

	// TODO things like CORS

	mux.Handle("/ping", httphelper.Log(pong()))

	todoDB := db.NewTodo(conn)
	mux.Handle("/todo/",
		httphelper.Log(
			http.StripPrefix("/todo",
				todos.New(todoDB).Routes(),
			)))
	return mux
}

func pong() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("PONG")); err != nil {
			slog.With("err", err).ErrorContext(r.Context(), "unable to respond to ping")
		}
	})
}
