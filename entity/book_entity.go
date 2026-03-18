package entity

import (
	"strings"

	"github.com/google/uuid"
)

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
	var names strings.Builder
	for _, author := range b.Authors {
		_, _ = names.WriteString(author.Name + ", ")
	}
	return names.String()
}

func (b *BookEntity) GetCategoryNames() string {
	var names strings.Builder
	for _, category := range b.Categories {
		_, _ = names.WriteString(category.Name + ", ")
	}
	return names.String()
}
