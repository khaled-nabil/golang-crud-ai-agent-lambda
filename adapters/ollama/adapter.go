package ollama

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
)

type (
	Ollama struct {
		model  string
		client *api.Client
	}
)

const (
	model                     = "nomic-embed-text"
	prefixEmbedSearchQuery    = "search_query: "
	prefixEmbedSearchDocument = "search_document: "
)

func New() (*Ollama, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("create Ollama client %w", err)
	}

	return &Ollama{
		model,
		client,
	}, nil
}

func (o *Ollama) EmbedSearchQuery(query string) ([]float64, error) {
	e, err := o.client.Embeddings(context.Background(), &api.EmbeddingRequest{
		Model:  o.model,
		Prompt: fmt.Sprintf("%s%s", prefixEmbedSearchQuery, query),
	})
	if err != nil {
		return nil, fmt.Errorf("embedding search query %w", err)
	}

	return e.Embedding, nil
}

func (o *Ollama) EmbedSearchDocument(text string) ([]float64, error) {
	e, err := o.client.Embeddings(context.Background(), &api.EmbeddingRequest{
		Model:  o.model,
		Prompt: fmt.Sprintf("%s%s", prefixEmbedSearchDocument, text),
	})
	if err != nil {
		return nil, fmt.Errorf("embedding search query %w", err)
	}

	return e.Embedding, nil
}
