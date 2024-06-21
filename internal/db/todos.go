package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/vankleefjim/go_experiment_jet/internal/db/.gen/things/public/model"
	. "github.com/vankleefjim/go_experiment_jet/internal/db/.gen/things/public/table"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/google/uuid"
)

type TodoDB struct {
	conn *sql.DB
}

func NewTodo(conn *sql.DB) *TodoDB { return &TodoDB{conn: conn} }

func (db *TodoDB) GetAll(ctx context.Context) ([]model.Todo, error) {
	todos := []model.Todo{}
	getAllStmt := Todo.SELECT(Todo.AllColumns)
	err := getAllStmt.QueryContext(ctx, db.conn, &todos)
	if err != nil {
		return nil, fmt.Errorf("unable to query all todos: %w", err)
	}
	return todos, nil
}

func (db *TodoDB) GetByID(ctx context.Context, id uuid.UUID) (model.Todo, error) {
	todos := []model.Todo{}
	getOneStmt := Todo.SELECT(Todo.AllColumns).WHERE(Todo.ID.EQ(UUID(id)))
	err := getOneStmt.QueryContext(ctx, db.conn, &todos)
	if err != nil {
		return model.Todo{}, fmt.Errorf("unable to query all todos: %w", err)
	}
	if len(todos) > 1 {
		return model.Todo{}, fmt.Errorf("found %d entries with ID %q", len(todos), id)
	}
	if len(todos) == 0 {
		return model.Todo{}, newTodoNotFoundError(fmt.Sprintf("id = %q", id))
	}
	return todos[0], nil
}

func (db *TodoDB) Create(ctx context.Context, todo model.Todo) error {
	insertStmt := Todo.INSERT(Todo.AllColumns).MODEL(todo)
	_, err := insertStmt.ExecContext(ctx, db.conn)
	if err != nil {
		return fmt.Errorf("unable to query all todos: %w", err)
	}
	// Inserting 1 new row: no need to check result

	return nil
}

func newTodoNotFoundError(identifier string) *NotFoundError {
	return &NotFoundError{
		Resource:   "todo",
		Identifier: identifier,
	}
}
