package db

import (
	"ai-agent/entity"
	"context"
	"fmt"

	"ai-agent/repositories/db/repo_dto"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

type (
	ChatRepository struct {
		agent *pgxpool.Pool
	}
)

const (
	chatTable            = "chat"
	chatQueryLimit int32 = 10
)

func NewChatRepository(db *PostgresRepo) *ChatRepository {
	return &ChatRepository{db.agent}
}

func (r *ChatRepository) StoreConversation(userID string, h *entity.ChatHistoryEntity, embedding []float32) (*entity.ChatHistoryEntity, error) {
	ctx := context.Background()

	var chat repo_dto.ChatDTO

	chatQuery := fmt.Sprintf(`
        INSERT INTO %s (user_id, message, response, embedding, created_at) 
        VALUES ($1, $2, $3, $4, NOW()) 
        RETURNING id, user_id, message, response, created_at
    `, chatTable)

	err := r.agent.
		QueryRow(ctx, chatQuery, userID, h.UserInput, h.Response, pgvector.NewVector(embedding)).
		Scan(&chat.ID, &chat.UserID, &chat.Message, &chat.Response, &chat.CreateAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert chat: %w", err)
	}

	return chat.ChatToHistoryContext(), nil
}

func (r *ChatRepository) GetUserHistory(id string) ([]entity.ChatHistoryEntity, error) {
	ctx := context.Background()

	query := fmt.Sprintf(`
	SELECT id, user_id, message, response, created_at FROM %s 
	WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2
	`, chatTable)

	var chats []repo_dto.ChatDTO

	q, err := r.agent.Query(ctx, query, id, chatQueryLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to query chat: %w", err)
	}
	defer q.Close()

	for q.Next() {
		var chat repo_dto.ChatDTO
		if err := q.Scan(
			&chat.ID,
			&chat.UserID,
			&chat.Message,
			&chat.Response,
			&chat.CreateAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan chat: %w", err)

		}

		chats = append(chats, chat)
	}

	return repo_dto.ChatListToHistoryContextList(chats), nil
}

func (r *ChatRepository) GetUserSimilarDocuments(userID string, embedding []float32) ([]entity.ChatHistoryEntity, error) {
	ctx := context.Background()

	vec := pgvector.NewVector(embedding)

	query := fmt.Sprintf(`
        SELECT 
            id
			user_id,
            message,
            response,
            (embedding <=> $2) AS vector_distance,
			created_at
        FROM %s 
        WHERE user_id = $1 
        AND vector_distance <= 0.5
        ORDER BY vector_distance ASC
        LIMIT $3`, chatTable)

	rows, err := r.agent.Query(ctx, query, userID, vec, chatQueryLimit)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}
	defer rows.Close()

	var history []repo_dto.ChatDTO
	for rows.Next() {
		var h repo_dto.ChatDTO
		var _vector_distance float32
		if err := rows.Scan(
			&h.ID,
			&h.UserID,
			&h.Message,
			&h.Response,
			&_vector_distance,
			&h.CreateAt,
		); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		history = append(history, h)
	}

	return repo_dto.ChatListToHistoryContextList(history), rows.Err()
}
