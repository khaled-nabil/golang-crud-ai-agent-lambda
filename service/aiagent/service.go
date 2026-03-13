package aiagent

import (
	"ai-agent/model/datamodels"
	"ai-agent/model/servicemodels"
	"fmt"
)

type (
	Service struct {
		agent datamodels.Gemini
		db    servicemodels.Persistence
	}
)

func New(agent datamodels.Gemini, db servicemodels.Persistence) *Service {
	return &Service{agent, db}
}

func (s *Service) SendMessageWithHistory(userID, message string) (string, error) {
	embeddedInput, err := s.agent.EmbedMessage(message)
	if err != nil {
		return "", fmt.Errorf("failed to embed user input: %w", err)
	}

	similarDocuments, err := s.db.GetUserSimilarDocuments(userID, embeddedInput)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve user similar documents: %w", err)
	}

	r, err := s.agent.Chat(message, datamodels.ChatListToHistoryContextList(similarDocuments))
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

	if err = s.db.StoreConversation(userID, r, embeddedConversation); err != nil {
		return "", err
	}

	return r.Response, err
}
