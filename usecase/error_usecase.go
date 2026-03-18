package usecase

import (
	"ai-agent/entity/errormodels"
	"net/http"
)

type (
	ErrorHandler struct {
	}
)

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (e *ErrorHandler) GetFormattedError(code errormodels.ErrorCodes, m ...string) errormodels.ErrCodeMapping {
	mapping, exists := errormodels.ErrMappings[code]
	if !exists {
		mapping = errormodels.ErrCodeMapping{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	if len(m) > 0 {
		mapping.Message = m[0]
	}

	return mapping
}
