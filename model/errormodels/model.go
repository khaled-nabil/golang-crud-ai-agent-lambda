package errormodels

import "net/http"

type (
	Errors interface {
		GetFormattedError(code ErrorCodes, debugMessage string) (int, FormattedError)
	}
	ErrorCodes string

	FormattedError struct {
		Code         ErrorCodes
		DebugMessage string
	}
)

const (
	ErrBadRequest ErrorCodes = "BAD_REQUEST"
	ErrNotFound   ErrorCodes = "NOT_FOUND"
	ErrGeneric    ErrorCodes = "GENERIC_ERROR"
)

var (
	ErrMappings = map[ErrorCodes]int{
		ErrBadRequest: http.StatusBadRequest,
		ErrNotFound:   http.StatusNotFound,
		ErrGeneric:    http.StatusInternalServerError,
	}
)
