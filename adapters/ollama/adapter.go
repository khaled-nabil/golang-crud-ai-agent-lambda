package ollama

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
)

type (
	OllamaAdapter struct {
		model  string
		client *api.Client
	}
)

const (
	model                     = "nomic-embed-text"
	prefixEmbedSearchQuery    = "search_query: "
	prefixEmbedSearchDocument = "search_document: "
)

func NewOllama() (*OllamaAdapter, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("create Ollama client %w", err)
	}
	
	return &OllamaAdapter{
		model,
		client,
	}, nil
}

func (o *OllamaAdapter) EmbedSearchQuery(query string) ([]float32, error) {
	e, err := o.client.Embeddings(context.Background(), &api.EmbeddingRequest{
		Model:  o.model,
		Prompt: fmt.Sprintf("%s%s", prefixEmbedSearchQuery, query),
	})
	if err != nil {
		return nil, fmt.Errorf("embedding search query %w", err)
	}

	return Float64To32(e.Embedding), nil
}

func (o *OllamaAdapter) EmbedSearchDocument(text string) ([]float32, error) {
	e, err := o.client.Embeddings(context.Background(), &api.EmbeddingRequest{
		Model:  o.model,
		Prompt: fmt.Sprintf("%s%s", prefixEmbedSearchDocument, text),
	})
	if err != nil {
		return nil, fmt.Errorf("embedding search query %w", err)
	}

	return Float64To32(e.Embedding), nil
}

func Float64To32(input []float64) []float32 {
    output := make([]float32, len(input))
    for i, v := range input {
        output[i] = float32(v)
    }
    return output
}