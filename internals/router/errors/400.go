package errors

import (
	"fmt"
	"net/http"
	"strings"
)

type ErrBadRequest struct {
	Msg string `json:"message"`
}

func NewErrBadRequest(format string, a ...string) RouteErrorHandler {
	return defaultErrorHandler{ErrBadRequest{Msg: fmt.Sprintf(format, a)}}
}
func (e ErrBadRequest) Error() string { return e.Msg }
func (e ErrBadRequest) Status() int   { return http.StatusBadRequest }

type ErrMissingParams struct {
	Params []string `json:"params"`
}

func NewErrMissingParams(params ...string) RouteErrorHandler {
	return defaultErrorHandler{ErrMissingParams{Params: params}}
}
func (e ErrMissingParams) Error() string {
	return fmt.Sprintf("Missing parameters: %s.", strings.Join(e.Params, ", "))
}
func (e ErrMissingParams) Status() int { return http.StatusBadRequest }

type ErrMethodNotAllowed struct {
	Method  string   `json:"method"`
	Allowed []string `json:"allowed"`
}

func NewErrMethodNotAllowed(method string, allowed ...string) RouteErrorHandler {
	return defaultErrorHandler{ErrMethodNotAllowed{Method: method, Allowed: allowed}}
}
func (e ErrMethodNotAllowed) Error() string {
	return fmt.Sprintf("Method %s not allowed. Allowed methods are: %s", e.Method, strings.Join(e.Allowed, ", "))
}
func (e ErrMethodNotAllowed) Status() int { return http.StatusMethodNotAllowed }
