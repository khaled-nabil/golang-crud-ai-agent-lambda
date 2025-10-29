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

var (
	userRole = "user"
	aiRole   = "ai"
)

func New(agent datamodels.Gemini, db servicemodels.AgentRepo) *Service {
	return &Service{agent, db}
}

func (s *Service) SendMessageWithHistory(userID, message string) (string, error) {
	if err := s.db.StoreChatMessage(userID, message, userRole); err != nil {
		return "", err
	}

	r, _, err := s.agent.Chat(message, nil)
	if err != nil {
		return "", err
	}

	if err = s.db.StoreChatMessage(userID, r, aiRole); err != nil {
		return "", err
	}

	return r, err
}
