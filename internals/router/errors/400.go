package errors

import (
	"fmt"
	"net/http"
	"strings"
)

type ErrMissingParams struct {
	defaultErr
	Params []string `json:"params"`
}

func NewErrMissingParams(params ...string) ErrMissingParams {
	return ErrMissingParams{Params: params}
}
func (e ErrMissingParams) Error() string {
	return fmt.Sprintf("Missing parameters: %s.", strings.Join(e.Params, ", "))
}
func (e ErrMissingParams) Status() int { return http.StatusBadRequest }
