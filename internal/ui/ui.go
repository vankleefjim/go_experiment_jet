package ui

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/vankleefjim/go_experiment_jet/internal/ui/components"
	"github.com/vankleefjim/go_experiment_jet/pkg/httphelper"
)

func RegisterRoutes(ctx context.Context, cfg Config, mux *http.ServeMux) *http.ServeMux {
	// Create the dependencies here.
	// If possible, use the ctx to control if something needs to be stopped or similar
	// Would be best if only the shared ones are here and the others
	// directly in the packages that define the routes.
	// dbConn := must(dbconn.SQLConnect(cfg.DB))
	// go func() {
	// 	<-ctx.Done()
	// 	cErr := dbConn.Close()
	// 	if cErr != nil {
	// 		slog.With("err", cErr).ErrorContext(ctx, "failed closing db connection")
	// 	}
	// }()

	// TODO things like CORS

	mux.Handle("/ping", httphelper.Log(pong()))

	// mux.Handle("/ui/",
	// 	httphelper.Log(
	// 		http.StripPrefix("/ui",
	// 			todos.New(todoDB).Routes(),
	// 		)))

	// TODO gzip this or does that already happen?
	mux.Handle("GET /ui", templ.Handler(components.Home()))
	return mux
}

func pong() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("PONG\n")); err != nil {
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
