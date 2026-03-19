package handler_dto

type (
	ChatRequest struct {
		Message string `json:"message"`
		UserID  string `json:"user_id"`
	}

	ChatResponse struct {
		Message string `json:"message"`
	}
)
