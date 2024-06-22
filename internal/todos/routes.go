package todos

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vankleefjim/go_experiment_jet/internal/db"
	"github.com/vankleefjim/go_experiment_jet/internal/httphelper"
	"github.com/vankleefjim/go_experiment_jet/pkg/collections"

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
			Get: httphelper.StructResponse[GetAllResponse](t.getAll),
			Put: httphelper.StructResponse[PutResponse](t.put),
		},
	))
	mux.HandleFunc("/{id}", httphelper.MethodPlexMiddleware(
		httphelper.MethodPlexer{
			Get: httphelper.StructResponse[GetOneResponse](t.get),
		},
	))

	return mux
}

func (t *TodoServer) getAll(r *http.Request) (*httphelper.OK[GetAllResponse], *httphelper.HTTPError) {
	ctx := r.Context()

	todos, err := t.db.GetAll(ctx)
	if err != nil {
		return nil, httphelper.NewError("unable to find todos", http.StatusInternalServerError, err)
	}

	return &httphelper.OK[GetAllResponse]{
		Body: GetAllResponse{
			Todos: collections.Map(todos, FromModel),
		}, Status: http.StatusOK}, nil
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

func (t *TodoServer) put(r *http.Request) (*httphelper.OK[PutResponse], *httphelper.HTTPError) {
	ctx := r.Context()

	todo := Todo{}
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		return nil, httphelper.NewError("invalid request body", http.StatusBadRequest, fmt.Errorf("unable to decode json body: %w", err))
	}

	err = todo.Validate()
	if err != nil {
		return nil, httphelper.NewError(err.Error(), http.StatusBadRequest, fmt.Errorf("validation failed: %w", err))
	}

	// Make sure to not accept ID from caller.
	todo.ID = uuid.New()

	err = t.db.Create(ctx, ToModel(todo))
	if err != nil {
		return nil, httphelper.NewError("unable to create todo", http.StatusInternalServerError, fmt.Errorf("unable to create todo: %w", err))
	}

	return &httphelper.OK[PutResponse]{
		Body:   PutResponse{Todo: todo},
		Status: http.StatusOK}, nil
}
