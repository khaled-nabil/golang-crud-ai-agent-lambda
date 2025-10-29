package servicemodels

import "ai-agent/model/datamodels"

type (
	AgentService interface {
		SendMessageWithHistory(userID, message string) (string, error)
	}

	AgentRepo interface {
		StoreConversation(userID string, history *datamodels.HistoryContext) error
		GetUserHistory(id string) ([]datamodels.HistoryContext, error)
	}

	RequestBody struct {
		Message string `json:"message"`
		UserID  string `json:"user_id"`
	}
)
