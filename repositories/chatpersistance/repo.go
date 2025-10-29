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
		UserID    string `json:"user_id" dynamodbav:"user_id"`
		Timestamp string `json:"timestamp" dynamodbav:"timestamp"`
		Message   string `json:"message" dynamodbav:"message"`
		Role      string `json:"role" dynamodbav:"role"`
	}
)

func New(pkg datamodels.DynamoDB) *Repo {
	return &Repo{pkg}
}

func (r *Repo) StoreChatMessage(userID, message, role string) error {
	chatMessage := ChatMessage{
		UserID:    userID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Message:   message,
		Role:      role,
	}

	err := r.pkg.StoreItem(chatMessage)

	return err
}
