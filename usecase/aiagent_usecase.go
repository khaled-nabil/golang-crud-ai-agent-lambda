package usecase

import (
	"ai-agent/interface"
	"fmt"
)

type (
	AIAgentUsecase struct {
		agent _interface.AgentAdapter
		cr    _interface.ChatRepo
	}
)

func NewAIAgentUsecase(agent _interface.AgentAdapter, cr _interface.ChatRepo) *AIAgentUsecase {
	return &AIAgentUsecase{agent, cr}
}

func (s *AIAgentUsecase) SendMessageWithHistory(userID, message, sysIns string) (string, error) {
	embeddedInput, err := s.agent.EmbedMessage(message)
	if err != nil {
		return "", fmt.Errorf("failed to embed user input: %w", err)
	}

	similarDocuments, err := s.cr.GetUserSimilarDocuments(userID, embeddedInput)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve user similar documents: %w", err)
	}

	r, err := s.agent.Chat(message, sysIns, similarDocuments)
	if err != nil {
		return "", fmt.Errorf("failed to send message to user: %w", err)
	}

	/*
	 * Current we embed user's input and AI response together
	 * TODO: assess if this should be improved, some alternatives
	 * - embed seperately, and later insert into two vector columns
	 * - use LLM to pull insights and important content from messages and insert as one
	 */
	embeddedConversation, err := s.agent.EmbedConverastion(r)
	if err != nil {
		return "", fmt.Errorf("failed to embed conversation: %w", err)
	}

	_, err = s.cr.StoreConversation(userID, r, embeddedConversation)
	if err != nil {
		return "", err
	}

	return r.Response, err
}
