package errors

import (
	"errors"
	"net/http"
)

type ErrInternal struct {
	Err string `json:"error"`
}

func NewErrInternal(err ...error) RouteErrorHandler {
	return defaultErrorHandler{ErrInternal{Err: errors.Join(err...).Error()}}
}
func (e ErrInternal) Error() string { return e.Err }
func (e ErrInternal) Status() int   { return http.StatusInternalServerError }
