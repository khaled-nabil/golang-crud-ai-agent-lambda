package db

import (
	"ai-agent/entity"
	"context"
	"fmt"
	"time"

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
	chatTable           = "chat"
	genaiEmbeddingTable = "documents_gemini"
	chatQueryLimit int32 = 10
)

func NewChatRepository(db *PostgresRepo) *ChatRepository {
	return &ChatRepository{db.agent}
}

func (r *ChatRepository) StoreConversation(userID string, h *entity.ChatHistoryEntity, embedding []float32) error {
	ctx := context.Background()

	tx, err := r.agent.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var chatID string
	var createdAt time.Time

	chatQuery := fmt.Sprintf(`
        INSERT INTO %s (user_id, message, response, created_at) 
        VALUES ($1, $2, $3, NOW()) 
        RETURNING id, created_at
    `, chatTable)

	err = tx.
		QueryRow(ctx, chatQuery, userID, h.UserInput, h.Response).
		Scan(&chatID, &createdAt)
	if err != nil {
		return fmt.Errorf("failed to insert chat: %w", err)
	}

	vec := pgvector.NewVector(embedding)

	embedQuery := fmt.Sprintf(`
        INSERT INTO %s (chat_id, user_id, embedding, created_at) 
        VALUES ($1, $2, $3, $4)
    `, genaiEmbeddingTable)

	_, err = tx.Exec(ctx, embedQuery, chatID, userID, vec, createdAt)
	if err != nil {
		return fmt.Errorf("failed to insert embedding: %w", err)
	}

	return tx.Commit(ctx)
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
            e.chat_id,
			e.user_id,
            c.message,
            c.response,
            e.created_at,
            (1 - (e.embedding <=> $2)) * 
            EXP(-EXTRACT(EPOCH FROM (NOW() - e.created_at))/86400 * 0.1) AS relevance_score
        FROM %s e
        JOIN %s c ON e.chat_id = c.id
        WHERE e.user_id = $1 
        AND 1 - (e.embedding <=> $2) >= 0.5
        ORDER BY relevance_score DESC
        LIMIT $3`, genaiEmbeddingTable, chatTable)

	rows, err := r.agent.Query(ctx, query, userID, vec, chatQueryLimit)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}
	defer rows.Close()

	var history []repo_dto.ChatDTO
	for rows.Next() {
		var h repo_dto.ChatDTO
		var _relevance_score float32
		if err := rows.Scan(
			&h.ID,
			&h.UserID,
			&h.Message,
			&h.Response,
			&h.CreateAt,
			&_relevance_score,
		); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		history = append(history, h)
	}

	return repo_dto.ChatListToHistoryContextList(history), rows.Err()
}
