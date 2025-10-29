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
