package errors

import (
	"ai-agent/model/errormodels"
	"net/http"
)

type (
	Errors struct {
		// implement logging interface later
	}
)

func New() *Errors {
	return &Errors{}
}

func GetFormattedError(code errormodels.ErrorCodes, debugMessage string) (int, errormodels.FormattedError) {
	s, ok := errormodels.ErrMappings[code]
	if !ok {
		s = http.StatusBadRequest
	}

	return s, errormodels.FormattedError{
		Code:         code,
		DebugMessage: debugMessage,
	}
}
