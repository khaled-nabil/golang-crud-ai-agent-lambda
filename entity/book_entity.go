package entity

import "github.com/google/uuid"

type (
	BookEntity struct {
		ID            *uuid.UUID
		Year          int16
		RatingCount   int
		PageCount     int
		AverageRating float32
		Title         string
		Subtitle      string
		Description   string
		Thumb         string
		Authors       []BookAuthorEntity
		Categories    []BookCategoryEntity
		Embedding     []float32
	}

	BookAuthorEntity struct {
		ID   *uuid.UUID
		Name string
	}

	BookCategoryEntity struct {
		ID   *uuid.UUID
		Name string
	}
)

func (b *BookEntity) GetAuthorNames() string {
	var names string
	for _, author := range b.Authors {
		names += author.Name + ", "
	}
	return names
}

func (b *BookEntity) GetCategoryNames() string {
	var names string
	for _, category := range b.Categories {
		names += category.Name + ", "
	}
	return names
}