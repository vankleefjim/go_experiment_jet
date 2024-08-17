package ui

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/vankleefjim/go_experiment_jet/internal/grpc/pbtodo"
	"github.com/vankleefjim/go_experiment_jet/internal/ui/components"
	"github.com/vankleefjim/go_experiment_jet/internal/ui/consts"
	"github.com/vankleefjim/go_experiment_jet/pkg/httphelper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ui struct {
	apiClient pbtodo.TodoServiceClient
}

func New(cfg Config) *ui {
	grpcClient, err := grpc.NewClient(cfg.APIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.With("err", err).Error("failed creating todo client")
		panic(err)
	}
	return &ui{
		apiClient: pbtodo.NewTodoServiceClient(grpcClient),
	}
}

func (u *ui) RegisterRoutes(mux *http.ServeMux) *http.ServeMux {

	// TODO things like CORS

	mux.Handle("/ping", httphelper.Log(pong()))

	// mux.Handle("/ui/",
	// 	httphelper.Log(
	// 		http.StripPrefix("/ui",
	// 			todos.New(todoDB).Routes(),
	// 		)))

	// TODO gzip this or does that already happen?
	mux.Handle("GET /ui", templ.Handler(components.Home()))

	mux.Handle("GET "+consts.PathTodos, u.getAllTodos())

	mux.Handle("PUT /ui/todo", httphelper.AddHeader("HX-Trigger", consts.TriggerTodosUpdate, templ.Handler(components.NYI())))
	return mux
}

func (u *ui) getAllTodos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		resp, err := u.apiClient.GetAll(ctx, &pbtodo.GetAllRequest{})
		if err != nil {
			slog.With("err", err).ErrorContext(ctx, "failed getting all todos")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		err = components.TodoList(resp.Todos).Render(ctx, w)
		if err != nil {
			slog.With("err", err, "todos", resp.Todos).ErrorContext(ctx, "failed rendering todolist")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}
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
