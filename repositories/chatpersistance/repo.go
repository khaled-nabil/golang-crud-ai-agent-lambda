package chatpersistance

import (
	"ai-agent/model/datamodels"
	"time"
)

type (
	Repo struct {
		pkg datamodels.DynamoDB
	}

	ChatMessage struct {
		UserID       string                     `json:"user_id" dynamodbav:"user_id"`
		Timestamp    int64                      `json:"timestamp" dynamodbav:"timestamp"`
		Conversation *datamodels.HistoryContext `json:"conversation" dynamodbav:"conversation"`
	}
)

var (
	limit int32 = 10
)

func New(pkg datamodels.DynamoDB) *Repo {
	return &Repo{pkg}
}

func (r *Repo) StoreConversation(id string, h *datamodels.HistoryContext) error {
	chatMessage := ChatMessage{
		UserID:       id,
		Timestamp:    time.Now().UTC().Unix(),
		Conversation: h,
	}

	err := r.pkg.StoreItem(chatMessage)

	return err
}

func (r *Repo) GetUserHistory(id string) ([]datamodels.HistoryContext, error) {
	items, err := r.pkg.RetrieveItems(id, limit)
	if err != nil {
		return nil, err
	}

	return items, nil
}
