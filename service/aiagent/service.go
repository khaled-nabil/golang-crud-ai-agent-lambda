package aiagent

import (
	"ai-agent/model/datamodels"
	"ai-agent/model/servicemodels"
)

type (
	Service struct {
		agent datamodels.Gemini
		db    servicemodels.AgentRepo
	}
)

func New(agent datamodels.Gemini, db servicemodels.AgentRepo) *Service {
	return &Service{agent, db}
}

func (s *Service) SendMessageWithHistory(userID, message string) (string, error) {
	h, err := s.agent.Chat(message, nil)
	if err != nil {
		return "", err
	}

	if err = s.db.StoreConversation(userID, h); err != nil {
		return "", err
	}

	return h.UserInput, err
}
