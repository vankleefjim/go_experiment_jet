package db

import "fmt"

var _ error = &NotFoundError{}

type NotFoundError struct {
	Resource   string
	Identifier string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("resource %q (%s) not found", e.Resource, e.Identifier)
}
