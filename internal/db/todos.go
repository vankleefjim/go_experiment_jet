package db

import (
	"context"
	"database/sql"
	"fmt"

	//. "github.com/go-jet/jet/v2/postgres"

	"jvk.com/things/internal/db/.gen/things/public/model"
	. "jvk.com/things/internal/db/.gen/things/public/table"
)

type TodosDB struct {
	conn *sql.DB
}

func NewTodos(conn *sql.DB) *TodosDB { return &TodosDB{conn: conn} }

func (db *TodosDB) GetAll(ctx context.Context) ([]model.Todos, error) {
	todos := []model.Todos{}
	getAllStmt := Todos.SELECT(Todos.AllColumns)
	err := getAllStmt.QueryContext(ctx, db.conn, &todos)
	if err != nil {
		return nil, fmt.Errorf("unable to query all todos: %w", err)
	}
	return todos, nil
}

func (db *TodosDB) Create(ctx context.Context, todo model.Todos) error {
	insertStmt := Todos.INSERT(Todos.AllColumns).MODEL(todo)
	_, err := insertStmt.ExecContext(ctx, db.conn)
	if err != nil {
		return fmt.Errorf("unable to query all todos: %w", err)
	}
	// Inserting 1 new row: no need to check result

	return nil
}
