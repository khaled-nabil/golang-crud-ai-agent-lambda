package geminipkg

import (
	"ai-agent/model/datamodels"
	"ai-agent/pkg/secretspkg"
	"context"
	"fmt"
	"log"

	"google.golang.org/genai"
)

type Gemini struct {
	model   string
	client  *genai.Client
	session *genai.Chat
	cfg     *secretspkg.AppConfig
}

var (
	temperature = float32(0.2)
	maxTokens   = int32(1024)
	system      = "You are a helpful AI assistant."
	ctx         = context.Background()
)

func New(cfg *secretspkg.AppConfig) (*Gemini, error) {
	apiKey := cfg.GeminiAPIKey
	modelName := cfg.ModelID

	log.Printf("Initializing Gemini client with model: %s", modelName)

	c, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return &Gemini{
		client: c,
		model:  modelName,
	}, nil
}

func (g *Gemini) createChat(history []*genai.Content) error {
	s, e := g.client.Chats.Create(ctx, g.model, &genai.GenerateContentConfig{
		Temperature:      &temperature,
		ResponseMIMEType: "application/json",
		MaxOutputTokens:  maxTokens,
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"message": {
					Type:        genai.TypeString,
					Description: "The AI-generated message response.",
				},
			},
		},
		SystemInstruction: &genai.Content{
			Role: "model",
			Parts: []*genai.Part{
				{
					Text: system,
				},
			},
		},
	}, history)

	if e != nil {
		return fmt.Errorf("failed to create chat session: %w", e)
	}

	g.session = s

	return nil
}

func (g *Gemini) Chat(userInput string, h []datamodels.HistoryContext) (string, []datamodels.HistoryContext, error) {
	if g.session == nil {
		if err := g.createChat(nil); err != nil {
			return "", nil, err
		}
	}

	uip := &genai.Part{
		Text: userInput,
	}

	resp, err := g.session.Send(ctx, uip)
	if err != nil {
		return "", nil, fmt.Errorf("failed to send message: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return "", nil, fmt.Errorf("no response received")
	}

	tr := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		tr += part.Text
	}

	h, err = extractHistory(&userInput, &tr, h)
	if err != nil {
		return "", nil, err
	}

	return tr, h, nil
}

func extractHistory(ui *string, r *string, h []datamodels.HistoryContext) ([]datamodels.HistoryContext, error) {
	if ui == nil || r == nil {
		return nil, fmt.Errorf("user input or response is empty")
	}

	h = append(h, datamodels.HistoryContext{
		UserInput: ui,
		Response:  r,
	})

	return h, nil
}
