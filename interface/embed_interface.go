package _interface

type (
	EmbeddingModel interface {
		EmbedSearchQuery(query string) ([]float64, error)
		EmbedSearchDocument(text string) ([]float64, error)
	}
)
