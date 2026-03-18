package _interface

import (
	"ai-agent/entity"
	"ai-agent/handler/handler_dto"
)

type (
	AgentAdapter interface {
		Chat(userInput string, systemPrompt string, h []entity.ChatHistoryEntity) (*entity.ChatHistoryEntity, error)
		EmbedMessage(text string) ([]float32, error)
		EmbedConverastion(h *entity.ChatHistoryEntity) ([]float32, error)
		RecommendBookFromList(userPrompt, systemPrompt string, b []entity.BookEntity) (*handler_dto.RecommendResponseDTO, error)
	}

	AgentUsecase interface {
		SendMessageWithHistory(userID, message, systemPrompt string) (string, error)
	}
)
