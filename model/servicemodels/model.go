package servicemodels

import "ai-agent/model/datamodels"

type (
	AgentService interface {
		SendMessageWithHistory(userID, message string) (string, error)
	}

	Persistence interface {
		StoreConversation(userID string, history *datamodels.HistoryContext, embedding []float32) error
		GetUserHistory(id string) ([]datamodels.Chat, error)
		GetUserSimilarDocuments(userID string, embedding []float32) ([]datamodels.Chat, error)
	}

	RequestBody struct {
		Message string `json:"message"`
		UserID  string `json:"user_id"`
	}
)
