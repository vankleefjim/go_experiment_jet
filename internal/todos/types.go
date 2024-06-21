package todos

import (
	"errors"
	"time"

	"github.com/vankleefjim/go_experiment_jet/internal/db/.gen/things/public/model"

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

func FromModel(in model.Todo) Todo {
	return Todo{
		ID:   in.ID,
		Task: in.Task,
		Due:  in.Due,
	}
}

func ToModel(in Todo) model.Todo {
	return model.Todo{
		ID:   in.ID,
		Task: in.Task,
		Due:  in.Due,
	}
}
