package gemini

import (
	"ai-agent/adapters/secrets"
	"ai-agent/entity"
	"bytes"
	"context"
	"fmt"

	"google.golang.org/genai"
)

type GeminiAdapter struct {
	model  string
	client *genai.Client
}

var (
	temperature   = float32(0.2)
	maxTokens     = int32(1024)
	embeddingSize = int32(1536)
)

const (
	embeddingModel             = "gemini-embedding-001"
	embeddingQeueryTaskType    = "RETRIEVAL_QUERY"
	embeddingRetrievalTaskType = "RETRIEVAL_DOCUMENT"
)

func NewGeminiAdapter(cfg *secrets.AppConfig) (*GeminiAdapter, error) {
	apiKey := cfg.GeminiAPIKey
	modelName := cfg.ModelID

	c, err := genai.NewClient(context.Background(), &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return &GeminiAdapter{
		client: c,
		model:  modelName,
	}, nil
}

func (g *GeminiAdapter) createChat(systemPrompt string, history []*genai.Content) (*genai.Chat, *context.Context, error) {
	ctx := context.Background()
	s, e := g.client.Chats.Create(ctx, g.model, &genai.GenerateContentConfig{
		Temperature:      &temperature,
		ResponseMIMEType: "text/plain",
		MaxOutputTokens:  maxTokens,
		SystemInstruction: &genai.Content{
			Role: genai.RoleModel,
			Parts: []*genai.Part{
				{
					Text: systemPrompt,
				},
			},
		},
	}, history)

	if e != nil {
		return nil, &ctx, fmt.Errorf("genai create chat session: %w", e)
	}

	return s, &ctx, nil
}

func (g *GeminiAdapter) Chat(userInput string, systemPrompt string, history []entity.ChatHistoryEntity) (*entity.ChatHistoryEntity, error) {
	if userInput == "" {
		return nil, fmt.Errorf("user input or response is empty")
	}

	h := transformHistoryToGeminiContent(history)

	c, ctx, err := g.createChat(systemPrompt, h)
	if err != nil {
		return nil, err
	}

	uip := &genai.Part{
		Text: userInput,
	}

	resp, err := c.Send(*ctx, uip)
	if err != nil {
		return nil, fmt.Errorf("genai send message: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("genai no response received")
	}

	var tr bytes.Buffer

	for _, part := range resp.Candidates[0].Content.Parts {
		_, err = tr.WriteString(part.Text)
		if err != nil {
			return nil, fmt.Errorf("genai write response: %w", err)
		}
	}

	return &entity.ChatHistoryEntity{
		UserInput: userInput,
		Response:  tr.String(),
	}, nil
}

func (g *GeminiAdapter) EmbedMessage(t string) ([]float32, error) {
	result, err := g.client.Models.EmbedContent(
		context.Background(),
		embeddingModel,
		[]*genai.Content{
			genai.NewContentFromText(t, genai.RoleUser),
		},
		&genai.EmbedContentConfig{
			OutputDimensionality: &embeddingSize,
			TaskType:             embeddingQeueryTaskType,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("genai embed message: %w", err)
	}

	embedding := result.Embeddings[0]

	return embedding.Values, nil
}

func (g *GeminiAdapter) EmbedConverastion(h *entity.ChatHistoryEntity) ([]float32, error) {
	result, err := g.client.Models.EmbedContent(
		context.Background(),
		embeddingModel,
		[]*genai.Content{
			genai.NewContentFromText(h.UserInput, genai.RoleUser),
			genai.NewContentFromText(h.Response, genai.RoleUser),
		},
		&genai.EmbedContentConfig{
			OutputDimensionality: &embeddingSize,
			TaskType:             embeddingRetrievalTaskType,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("genai embed conversation: %w", err)
	}

	embedding := result.Embeddings[0]

	return embedding.Values, nil
}

func transformHistoryToGeminiContent(h []entity.ChatHistoryEntity) []*genai.Content {
	var history []*genai.Content

	for _, item := range h {
		if item.UserInput != "" {
			history = append(history, genai.NewContentFromText(item.UserInput, genai.RoleUser))
		}
		if item.Response != "" {
			history = append(history, genai.NewContentFromText(item.Response, genai.RoleModel))
		}
	}
	return history
}
