package aiagentmodel

type (
	AgentService interface {
		SendMessageWithHistory(message string) (string, error)
	}

	RequestBody struct {
		Message string `json:"message"`
	}
)
