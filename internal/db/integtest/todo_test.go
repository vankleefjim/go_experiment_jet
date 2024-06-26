package integtest

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/vankleefjim/go_experiment_jet/internal/db/.gen/things/public/model"
	"github.com/vankleefjim/go_experiment_jet/pkg/thelper"
)

func Test_Todo_CreateAndGet(t *testing.T) {
	ctx := context.Background()
	t1 := &model.Todo{
		Task: "test something",
		Due:  ptr(time.Now()),
	}
	if err := todoDB.Create(ctx, t1); err != nil {
		t.Fatal(err)
	}

	if t1.ID == uuid.Nil {
		t.Error("expected ID to be nonnil")
	}

	t1R, err := todoDB.GetByID(ctx, t1.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(t1, t1R, thelper.FixMonotonicTimePtr()); diff != "" {
		t.Errorf("response. Diff(+got-want):\n%s", diff)
	}
}

func ptr[T any](v T) *T { return &v }
