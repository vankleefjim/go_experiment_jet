package api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/vankleefjim/go_experiment_jet/internal/db"
	"github.com/vankleefjim/go_experiment_jet/internal/grpc/pbtodo"
	"github.com/vankleefjim/go_experiment_jet/internal/todos"
	"github.com/vankleefjim/go_experiment_jet/pkg/dbconn"
	"github.com/vankleefjim/go_experiment_jet/pkg/httphelper"
	"google.golang.org/grpc"

	"log/slog"
)

type API struct {
	dbConn     *sql.DB
	todoServer *todos.TodoServer
}

func New(cfg Config) *API {
	// Create the dependencies here.
	// Would be best if only the shared ones are here and the others
	// directly in the packages that define the routes.
	dbConn := must(dbconn.SQLConnect(cfg.DB))
	return &API{
		dbConn:     dbConn,
		todoServer: todos.New(db.NewTodo(dbConn)),
	}
}

func (a *API) Shutdown(_ context.Context) {
	cErr := a.dbConn.Close()
	if cErr != nil {
		slog.With("err", cErr).Error("failed closing db connection")
	}
}

func (a *API) RegisterRoutes(mux *http.ServeMux) *http.ServeMux {
	runtime.NewServeMux()

	// TODO things like CORS
	mux.Handle("/ping", httphelper.Log(pong()))

	mux.Handle("/todo/",
		httphelper.Log(
			http.StripPrefix("/todo",
				a.todoServer.Routes(),
			)))
	return mux
}

func (a *API) RegisterGRPC(grpcServer *grpc.Server) {
	pbtodo.RegisterTodoServiceServer(grpcServer, a.todoServer)
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
