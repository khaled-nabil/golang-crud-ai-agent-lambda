package db

import (
	"ai-agent/model/datamodels"
	"fmt"
)

var (
	limit int32 = 10
)

const (
	chatTable = "chat"
)

func (r *Repository) StoreConversation(id string, h *datamodels.HistoryContext) error {
	query := `INSERT INTO $1 (user_id, message, response) VALUES ($2, $3, $4)`

	_, err := r.agent.Exec(r.ctx, query, chatTable, id, h.UserInput, h.Response)

	return fmt.Errorf("failed to store conversation: %w", err)
}

func (r *Repository) GetUserHistory(id string) ([]datamodels.Chat, error) {
	query := `SELECT id, user_id, message, response, created_at FROM $1 WHERE user_id = $2 ORDER BY created_at DESC LIMIT $3`

	var chats []datamodels.Chat

	q, err := r.agent.Query(r.ctx, query, chatTable, id, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query chat: %w", err)
	}
	defer q.Close()

	for q.Next() {
		var chat datamodels.Chat
		if err := q.Scan(
			chat.ID,
			chat.UserID,
			chat.Message,
			chat.Response,
			chat.CreateAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan chat: %w", err)

		}

		chats = append(chats, chat)
	}

	return chats, nil
}
