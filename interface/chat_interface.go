package _interface

import (
	"ai-agent/entity"
)

type (
	ChatRepo interface {
		StoreConversation(userID string, history *entity.ChatHistoryEntity, embedding []float32) error
		GetUserHistory(id string) ([]entity.ChatHistoryEntity, error)
		GetUserSimilarDocuments(userID string, embedding []float32) ([]entity.ChatHistoryEntity, error)
	}
)
