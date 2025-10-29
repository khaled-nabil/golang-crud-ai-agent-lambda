package servicemodels

type (
	AgentService interface {
		SendMessageWithHistory(userID, message string) (string, error)
	}

	AgentRepo interface {
		StoreChatMessage(userID, message, role string) error
	}

	RequestBody struct {
		Message string `json:"message"`
		UserID  string `json:"user_id"`
	}
)
