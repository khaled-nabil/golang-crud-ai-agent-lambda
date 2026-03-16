package repo_dto

import (
	"ai-agent/entity"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type (
	ChatDTO struct {
		ID       *pgtype.UUID `json:"id,omitempty"`
		UserID   pgtype.UUID  `json:"user_id"`
		Message  string       `json:"message"`
		Response string       `json:"response"`
		CreateAt *time.Time   `json:"created_at,omitempty"`
	}
)

func (c *ChatDTO) ChatToHistoryContext() *entity.ChatHistoryEntity {
	return &entity.ChatHistoryEntity{
		UserInput: c.Message,
		Response:  c.Response,
	}
}

func ChatListToHistoryContextList(c []ChatDTO) []entity.ChatHistoryEntity {
	var history []entity.ChatHistoryEntity

	for _, chat := range c {
		history = append(history, *chat.ChatToHistoryContext())
	}

	return history
}
