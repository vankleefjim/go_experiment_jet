package api

import (
	"context"
	"net/http"

	"github.com/vankleefjim/go_experiment_jet/internal/db"
	"github.com/vankleefjim/go_experiment_jet/internal/todos"
	"github.com/vankleefjim/go_experiment_jet/pkg/dbconn"
	"github.com/vankleefjim/go_experiment_jet/pkg/httphelper"

	"log/slog"
)

func RegisterRoutes(ctx context.Context, cfg Config, mux *http.ServeMux) *http.ServeMux {
	// Create the dependencies here.
	// If possible, use the ctx to control if something needs to be stopped or similar
	// Would be best if only the shared ones are here and the others
	// directly in the packages that define the routes.
	dbConn := must(dbconn.SQLConnect(cfg.DB))
	go func() {
		<-ctx.Done()
		cErr := dbConn.Close()
		if cErr != nil {
			slog.With("err", cErr).ErrorContext(ctx, "failed closing db connection")
		}
	}()

	// TODO things like CORS

	mux.Handle("/ping", httphelper.Log(pong()))

	todoDB := db.NewTodo(dbConn)
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

func must[T any](x T, err error) T {
	if err != nil {
		panic(err)
	}
	return x
}
