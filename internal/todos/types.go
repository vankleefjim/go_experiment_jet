package todos

import (
	"errors"
	"time"

	"github.com/vankleefjim/go_experiment_jet/internal/db/.gen/things/public/model"
	"github.com/vankleefjim/go_experiment_jet/internal/grpc/pbtodo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
)

type Todo struct {
	ID   uuid.UUID  `json:"id"` // Cannot be set by caller.
	Task string     `json:"task"`
	Due  *time.Time `json:"due"`
}

func (t Todo) Validate() error {
	errs := []error{}

	if t.Task == "" {
		errs = append(errs, errors.New("task may not be empty"))
	}
	if t.Due != nil && t.Due.Before(time.Now()) {
		errs = append(errs, errors.New("due may only be in the future"))
	}

	return errors.Join(errs...)
}

type Todos []Todo

type GetAllResponse struct {
	Todos Todos `json:"todos"`
}

type PutResponse struct {
	Todo Todo `json:"todo"`
}

type GetOneResponse struct {
	Todo Todo `json:"todo"`
}

func FromModel(in *model.Todo) Todo {
	return Todo{
		ID:   in.ID,
		Task: in.Task,
		Due:  in.Due,
	}
}

func PBID(in uuid.UUID) *pbtodo.UUID {
	return &pbtodo.UUID{Value: in.String()}
}

func PBTime(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func PBFromModel(in *model.Todo) *pbtodo.Todo {
	return &pbtodo.Todo{
		ID:   PBID(in.ID),
		Task: in.Task,
		Due:  PBTime(in.Due),
	}
}

func ToModel(in Todo) *model.Todo {
	return &model.Todo{
		ID:   in.ID,
		Task: in.Task,
		Due:  in.Due,
	}
}
