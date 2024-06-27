package integtest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/vankleefjim/go_experiment_jet/internal/db"
	"github.com/vankleefjim/go_experiment_jet/internal/db/.gen/things/public/model"
	"github.com/vankleefjim/go_experiment_jet/pkg/thelper"
)

func Test_Todo_CreateGetDelete(t *testing.T) {
	ctx := context.Background()
	t1 := &model.Todo{
		ID:   uuid.New(),
		Task: "test something",
		Due:  ptr(time.Now()),
	}
	if err := todoDB.Create(ctx, t1); err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := todoDB.Delete(ctx, t1.ID)
		if err != nil {
			t.Error(err)
		}
	}()

	t1R, err := todoDB.GetByID(ctx, t1.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(t1, t1R, thelper.FixMonotonicTimePtr()); diff != "" {
		t.Errorf("response. Diff(+got-want):\n%s", diff)
	}
}

func Test_Todo_GetAll(t *testing.T) {
	ctx := context.Background()
	t1 := &model.Todo{
		ID:   uuid.New(),
		Task: "test something",
		Due:  ptr(time.Now()),
	}
	if err := todoDB.Create(ctx, t1); err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := todoDB.Delete(ctx, t1.ID)
		if err != nil {
			t.Error(err)
		}
	}()
	t2 := &model.Todo{
		ID:   uuid.New(),
		Task: "test something else",
		Due:  ptr(time.Now()),
	}
	if err := todoDB.Create(ctx, t2); err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := todoDB.Delete(ctx, t2.ID)
		if err != nil {
			t.Error(err)
		}
	}()

	allTodos, err := todoDB.GetAll(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff([]*model.Todo{t1, t2}, allTodos,
		thelper.FixMonotonicTimePtr(), cmpopts.SortSlices(todoIDLess)); diff != "" {
		t.Errorf("response. Diff(+got-want):\n%s", diff)
	}
}

func Test_Todo_DeleteNonExistent(t *testing.T) {
	ctx := context.Background()
	t1 := &model.Todo{
		ID:   uuid.New(),
		Task: "test something",
		Due:  ptr(time.Now()),
	}
	if err := todoDB.Create(ctx, t1); err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := todoDB.Delete(ctx, t1.ID)
		if err != nil {
			t.Error(err)
		}
	}()

	err := todoDB.Delete(ctx, uuid.New())
	if !errors.As(err, ptr(&db.NotFoundError{})) {
		t.Errorf("expected err of type %T but got %v", &db.NotFoundError{}, err)
	}

	// to check it wasn't deleted
	t1R, err := todoDB.GetByID(ctx, t1.ID)
	if err != nil {
		t.Fatal(err)
	}
	if t1R == nil {
		t.Error("found todo unexpectedly nil")
	}
}

// helper funcs

func todoIDLess(x, y *model.Todo) bool {
	return x.ID.String() < y.ID.String()
}

func ptr[T any](v T) *T { return &v }
