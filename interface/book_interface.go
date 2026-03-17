package _interface

import (
	"ai-agent/entity"
)

type (
	BookRepo interface {
		InsertBookIfNotExists(book *entity.BookEntity) error
	}

	BookUsecase interface {
		Insert(book *entity.BookEntity) error
	}
)
