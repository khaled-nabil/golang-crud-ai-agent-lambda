package _interface

type (
	EmbeddingModel interface {
		EmbedSearchQuery(query string) ([]float32, error)
		EmbedSearchDocument(text string) ([]float32, error)
	}
)
