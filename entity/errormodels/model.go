package errormodels

import "net/http"

type (
	Errors interface {
		GetFormattedError(code ErrorCodes, msg ...string) ErrCodeMapping
	}
	ErrorCodes string

	ErrCodeMapping struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

const (
	ErrBadRequest     ErrorCodes = "BAD_REQUEST"
	ErrNotFound       ErrorCodes = "NOT_FOUND"
	ErrGeneric        ErrorCodes = "GENERIC_ERROR"
	ErrUnauthorized   ErrorCodes = "UNAUTHORIZED"
)

var (
	ErrMappings = map[ErrorCodes]ErrCodeMapping{
		ErrBadRequest:   {Code: http.StatusBadRequest, Message: "Bad Request"},
		ErrNotFound:     {Code: http.StatusNotFound, Message: "Not Found"},
		ErrGeneric:      {Code: http.StatusInternalServerError, Message: "Internal Server Error"},
		ErrUnauthorized: {Code: http.StatusUnauthorized, Message: "user not authenticated or unauthorized"},
	}
)
