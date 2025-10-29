package aiagent

import "ai-agent/pkg/geminipkg"

type (
	Service struct {
		agent *geminipkg.Gemini
	}
)

func New(agent *geminipkg.Gemini) *Service {
	return &Service{agent: agent}
}

func (s *Service) SendMessageWithHistory(message string) (string, error) {
	r, _, err := s.agent.Chat(message, nil)

	return r, err
}
