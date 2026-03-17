package _interface

import (
	"ai-agent/entity"
)

type (
	BookRepo interface {
		InsertBookIfNotExists(book *entity.BookEntity) error
		SearchForRelevantBook(embedding []float32) ([]entity.BookEntity, error)
	}

	BookUsecase interface {
		Insert(book *entity.BookEntity) error
		GetBookRecommendations(prompt string) ([]entity.BookEntity, error)
	}
)
