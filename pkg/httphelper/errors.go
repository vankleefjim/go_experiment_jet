package httphelper

import "fmt"

var _ error = &HTTPError{}

func NewError(message string, code int, cause error) *HTTPError {
	return &HTTPError{Message: message, Code: code, Cause: cause}
}

type HTTPError struct {
	Message string
	Code    int
	Cause   error // Not for external consumption.
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTPStatusCode [%d]: %s. \n Cause: %s", e.Code, e.Message, e.Cause)
}
