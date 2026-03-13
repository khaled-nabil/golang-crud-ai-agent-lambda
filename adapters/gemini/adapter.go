package gemini

import (
	"ai-agent/model/datamodels"
	"context"
	"fmt"

	"google.golang.org/genai"
)

type Gemini struct {
	model  string
	client *genai.Client
}

var (
	temperature   = float32(0.2)
	maxTokens     = int32(1024)
	system        = "You are a helpful AI assistant. Use the history context when possible to answer questions."
	embeddingSize = int32(1536)
)

const (
	emddingModel = "gemini-embedding-001"
)

func New(cfg *datamodels.AppConfig) (*Gemini, error) {
	apiKey := cfg.GeminiAPIKey
	modelName := cfg.ModelID

	c, err := genai.NewClient(context.Background(), &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return &Gemini{
		client: c,
		model:  modelName,
	}, nil
}

func (g *Gemini) createChat(history []*genai.Content) (*genai.Chat, *context.Context, error) {
	ctx := context.Background()
	s, e := g.client.Chats.Create(ctx, g.model, &genai.GenerateContentConfig{
		Temperature:      &temperature,
		ResponseMIMEType: "text/plain",
		MaxOutputTokens:  maxTokens,
		SystemInstruction: &genai.Content{
			Role: genai.RoleModel,
			Parts: []*genai.Part{
				{
					Text: system,
				},
			},
		},
	}, history)

	if e != nil {
		return nil, &ctx, fmt.Errorf("failed to create chat session: %w", e)
	}

	return s, &ctx, nil
}

func (g *Gemini) Chat(userInput string, history []datamodels.HistoryContext) (*datamodels.HistoryContext, error) {
	if userInput == "" {
		return nil, fmt.Errorf("user input or response is empty")
	}

	h := transformHistoryToGeminiContent(history)

	c, ctx, err := g.createChat(h)
	if err != nil {
		return nil, err
	}

	uip := &genai.Part{
		Text: userInput,
	}

	resp, err := c.Send(*ctx, uip)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response received")
	}

	tr := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		tr += part.Text
	}

	return &datamodels.HistoryContext{
		UserInput: userInput,
		Response:  tr,
	}, nil
}

func (g *Gemini) EmbedMessage(t string) ([]float32, error) {
	result, err := g.client.Models.EmbedContent(
		context.Background(),
		emddingModel,
		[]*genai.Content{
			genai.NewContentFromText(t, genai.RoleUser),
		},
		&genai.EmbedContentConfig{OutputDimensionality: &embeddingSize},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to embed message: %w", err)
	}

	embedding := result.Embeddings[0]

	return embedding.Values, nil
}

func transformHistoryToGeminiContent(h []datamodels.HistoryContext) []*genai.Content {
	var history []*genai.Content

	for _, item := range h {
		if item.UserInput != "" {
			history = append(history, &genai.Content{
				Role: genai.RoleUser,
				Parts: []*genai.Part{
					{
						Text: item.UserInput,
					},
				},
			})
		}
		if item.Response != "" {
			history = append(history, &genai.Content{
				Role: genai.RoleModel,
				Parts: []*genai.Part{
					{
						Text: item.Response,
					},
				},
			})
		}
	}
	return history
}
