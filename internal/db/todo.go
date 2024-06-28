package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/vankleefjim/go_experiment_jet/internal/db/.gen/things/public/model"
	. "github.com/vankleefjim/go_experiment_jet/internal/db/.gen/things/public/table"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"

	"github.com/google/uuid"
)

type TodoDB struct {
	conn *sql.DB
}

func NewTodo(conn *sql.DB) *TodoDB { return &TodoDB{conn: conn} }

func (db *TodoDB) GetAll(ctx context.Context) ([]*model.Todo, error) {
	// TODO:
	// This loads everything into memory at once, which is likely not very smart.
	// I didn't find jet to natively support a cursor or something.
	todos := []*model.Todo{}
	getAllStmt := Todo.SELECT(Todo.AllColumns)
	err := getAllStmt.QueryContext(ctx, db.conn, &todos)
	if err != nil {
		return nil, fmt.Errorf("unable to query all todos: %w", err)
	}
	return todos, nil
}

func (db *TodoDB) GetByID(ctx context.Context, id uuid.UUID) (*model.Todo, error) {
	return db.getOne(ctx, id, db.conn)
}

func (db *TodoDB) Create(ctx context.Context, todo *model.Todo) error {
	if todo.ID == uuid.Nil {
		return errors.New("ID must not be empty")
	}
	insertStmt := Todo.INSERT(Todo.AllColumns).MODEL(todo)
	_, err := insertStmt.ExecContext(ctx, db.conn)
	if err != nil {
		return fmt.Errorf("unable to query all todos: %w", err)
	}
	// Inserting 1 new row: no need to check result

	return nil
}

func (db *TodoDB) Delete(ctx context.Context, id uuid.UUID) error {
	deleteStmt := Todo.DELETE().WHERE(Todo.ID.EQ(UUID(id)))
	result, err := deleteStmt.ExecContext(ctx, db.conn)
	if err != nil {
		return fmt.Errorf("unable to delete todo with id %q: %w", id, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows affected for deletion of %q: %w", id, err)
	}
	if rowsAffected != 1 {
		return newTodoNotFoundError("id = " + id.String())
	}
	return nil
}

func (db *TodoDB) SetDue(ctx context.Context, id uuid.UUID, due time.Time) (err error) {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to begin tx: %w", err)
	}
	defer func() {
		// TODO make this a common method?
		if err != nil {
			rErr := tx.Rollback()
			if rErr != nil {
				err = errors.Join(err, fmt.Errorf("unable to rollback: %w", rErr))
			}
			return
		}
	}()

	_, err = db.getOne(ctx, id, tx)
	if err != nil {
		return err
	}

	updateDueStmt := Todo.UPDATE().
		SET(Todo.Due.SET(TimestampzT(due))).
		WHERE(Todo.ID.EQ(UUID(id)))
	result, err := updateDueStmt.ExecContext(ctx, tx)
	if err != nil {
		return fmt.Errorf("unable to update %q: %w", id, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows affected for %q: %w", id, err)
	}
	if rowsAffected != 1 {
		return fmt.Errorf("expected 1 affected row but had %d for %q", rowsAffected, id)
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("unable to commit transaction: %w", err)
	}
	return nil
}

func (db *TodoDB) getOne(ctx context.Context, id uuid.UUID, conn qrm.Queryable) (*model.Todo, error) {
	todo := &model.Todo{}

	getByIDStmt := Todo.SELECT(Todo.AllColumns).WHERE(Todo.ID.EQ(UUID(id)))
	err := getByIDStmt.QueryContext(ctx, conn, todo)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, newTodoNotFoundError(fmt.Sprintf("id = %q", id))
		}
		return nil, fmt.Errorf("unable to query all todos: %w", err)
	}
	return todo, nil
}

func newTodoNotFoundError(identifier string) *NotFoundError {
	return &NotFoundError{
		Resource:   "todo",
		Identifier: identifier,
	}
}
