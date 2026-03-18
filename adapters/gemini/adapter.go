package gemini

import (
	"ai-agent/adapters/secrets"
	"ai-agent/entity"
	"ai-agent/handler/handler_dto"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

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

func (g *GeminiAdapter) RecommendBookFromList(userPrompt, systemPrompt string, books []entity.BookEntity) (*handler_dto.RecommendResponseDTO, error) {
	ctx := context.Background()
	minBooks := int64(1)
	maxBooks := int64(5)

	c, err := g.client.Chats.Create(ctx, g.model, &genai.GenerateContentConfig{
		Temperature:      &temperature,
		ResponseMIMEType: "application/json",
		MaxOutputTokens:  maxTokens,
		SystemInstruction: &genai.Content{
			Role: genai.RoleModel,
			Parts: []*genai.Part{
				{
					Text: systemPrompt,
				},
			},
		},
		ResponseJsonSchema: &genai.Schema{
			Type: "object",
			Properties: map[string]*genai.Schema{
				"description": {
					Type:        "string",
					Description: "Description of the set of books selected and the user's prompt",
				},
				"recommended_books": {
					Type:     "array",
					MaxItems: &maxBooks,
					MinItems: &minBooks,
					Items: &genai.Schema{
						Type: "object",
						Properties: map[string]*genai.Schema{
							"title": {
								Type:        "string",
								Description: "title of the book",
							},
							"author": {
								Type:        "string",
								Description: "author of the book",
							},
							"genre": {
								Type:        "string",
								Description: "genre of the book",
							},
							"summary": {
								Type:        "string",
								Description: "summary of the book based on the description",
							},
							"reason": {
								Type:        "string",
								Description: "The reason this book was selected and how it matches the user's expectations",
							},
						},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("genai create chat session: %w", err)
	}

	var contextBuilder strings.Builder
	_, err = contextBuilder.WriteString("REFERENCE DOCUMENTS:\n")
	if err != nil {
		return nil, fmt.Errorf("genai write reference documents: %w", err)
	}
	for _, b := range books {
		_, err = contextBuilder.WriteString(transformBookToPrompt(&b))
		if err != nil {
			return nil, fmt.Errorf("genai transform book to prompt: %w", err)
		}
	}

	finalPrompt := fmt.Sprintf(
		"%s\n\nUSER QUESTION: %s\n\nREMINDER: Use ONLY the books listed above for your suggestions.",
		contextBuilder.String(),
		userPrompt,
	)

	resp, err := c.Send(ctx, genai.NewPartFromText(finalPrompt))
	if err != nil {
		return nil, fmt.Errorf("genai send message: %w", err)
	}

	var rawJSON bytes.Buffer

	for _, part := range resp.Candidates[0].Content.Parts {
		_, err = rawJSON.WriteString(part.Text)
		if err != nil {
			return nil, fmt.Errorf("genai write response: %w", err)
		}
	}

	if rawJSON.String() == "" {
		return nil, fmt.Errorf("empty response from model")
	}

	var recommendation handler_dto.RecommendResponseDTO
	if err := json.Unmarshal(rawJSON.Bytes(), &recommendation); err != nil {
		return nil, fmt.Errorf("parsing model response into DTO: %w", err)
	}

	return &recommendation, nil
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

func transformBookToPrompt(b *entity.BookEntity) string {
	return fmt.Sprintf("--- BOOK ID: %d ---\nTitle: %s\nAuthor: %s\nGenre: %s\nSummary: %s\n\n", *b.ID, b.Title, b.GetAuthorNames(), b.GetCategoryNames(), b.Description)
}
