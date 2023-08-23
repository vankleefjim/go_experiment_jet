package todos

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"jvk.com/things/internal/db/.gen/things/public/model"
)

type Todo struct {
	ID   uuid.UUID  `json:"id"`
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

func FromModel(in model.Todos) Todo {
	return Todo{
		ID:   in.ID,
		Task: in.Task,
		Due:  in.Due,
	}
}

func ToModel(in Todo) model.Todos {
	return model.Todos{
		ID:   in.ID,
		Task: in.Task,
		Due:  in.Due,
	}
}
