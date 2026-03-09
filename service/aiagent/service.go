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
	h, err := s.db.GetUserHistory(userID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve user history: %w", err)
	}

	r, err := s.agent.Chat(message, datamodels.ChatListToHistoryContextList(h))
	if err != nil {
		return "", fmt.Errorf("failed to send message to user: %w", err)
	}

	e, err := s.agent.EmbedMessage(message, r.Response)
	if err != nil {
		return "", fmt.Errorf("failed to embed message: %w", err)
	}

	if err = s.db.StoreConversation(userID, r, e); err != nil {
		return "", err
	}

	return r.Response, err
}
