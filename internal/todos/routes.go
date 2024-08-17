package todos

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/vankleefjim/go_experiment_jet/internal/db"
	"github.com/vankleefjim/go_experiment_jet/internal/grpc/pbtodo"
	"github.com/vankleefjim/go_experiment_jet/pkg/collections"
	"github.com/vankleefjim/go_experiment_jet/pkg/httphelper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

type TodoServer struct {
	db *db.TodoDB
}

func New(db *db.TodoDB) *TodoServer { return &TodoServer{db: db} }

func (t *TodoServer) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", httphelper.MethodPlexMiddleware(
		httphelper.MethodPlexer{
			//Get: httphelper.StructResponse[GetAllResponse](t.getAll),
			//Put: httphelper.StructResponse[PutResponse](t.put),
		},
	))
	mux.HandleFunc("/{id}", httphelper.MethodPlexMiddleware(
		httphelper.MethodPlexer{
			Get: httphelper.StructResponse[GetOneResponse](t.get),
		},
	))

	return mux
}

func (t *TodoServer) GetAll(ctx context.Context, r *pbtodo.GetAllRequest) (*pbtodo.GetAllResponse, error) {
	todos, err := t.db.GetAll(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to find todos")
	}

	return &pbtodo.GetAllResponse{
		Todos: collections.Map(todos, PBFromModel),
	}, nil
}

func (t *TodoServer) get(r *http.Request) (*httphelper.OK[GetOneResponse], *httphelper.HTTPError) {
	ctx := r.Context()

	idS := r.PathValue("id")
	id, err := uuid.Parse(idS)
	if err != nil {
		return nil, httphelper.NewError("invalid id: "+idS, http.StatusBadRequest, fmt.Errorf("unable to parse id %q: %w", idS, err))
	}

	todo, err := t.db.GetByID(ctx, id)
	if err != nil {
		return nil, httphelper.NewError("unable to find todos", http.StatusInternalServerError, err)
	}

	return &httphelper.OK[GetOneResponse]{
		Body: GetOneResponse{
			Todo: FromModel(todo),
		},
		Status: http.StatusOK}, nil
}

func (t *TodoServer) Put(ctx context.Context, r *pbtodo.PutRequest) (*pbtodo.PutResponse, error) {
	err := Validate(r.Todo)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Errorf("validation failed: %w", err).Error())
	}

	// Make sure to not accept ID from caller.
	newID := uuid.New()
	r.Todo.ID = PBID(newID)
	err = t.db.Create(ctx, ToModel(r.Todo, newID))
	if err != nil {
		slog.With("err", err, "todo", r.Todo).ErrorContext(ctx, "unable to create todo")
		return nil, status.Error(codes.Internal, "unable to create todo")
	}

	return &pbtodo.PutResponse{Todo: r.Todo}, nil
}
