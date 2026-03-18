package repo_dto

import (
	"ai-agent/entity"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type (
	BookDTO struct {
		ID          *pgtype.UUID `json:"id,omitempty"`
		Title       string       `json:"title"`
		Subtitle    string       `json:"subtitle,omitempty"`
		Description string       `json:"description"`
		Embedding   []float32    `json:"embedding"`
		Thumb       string       `json:"thumbnail"`
		PubYear     int8         `json:"published_year"`
		Rating      float32      `json:"average_rating"`
		RatingCount int          `json:"rating_count"`
		PageCount   int16        `json:"num_pages"`
		CreatedAt   *time.Time   `json:"created_at,omitempty"`
	}

	BookAuthorDTO struct {
		ID   *pgtype.UUID `json:"id,omitempty"`
		Name string       `json:"name"`
	}

	BookCategoryDTO struct {
		ID   *pgtype.UUID `json:"id,omitempty"`
		Name string       `json:"name"`
	}

	BookAuthorMapDTO struct {
		BookID   *pgtype.UUID `json:"book_id"`
		AuthorID *pgtype.UUID `json:"author_id"`
	}

	BookCategoryMapDTO struct {
		BookID     *pgtype.UUID `json:"book_id"`
		CategoryID *pgtype.UUID `json:"category_id"`
	}
)

func (b *BookDTO) ToBookEntity(authors []BookAuthorDTO, categories []BookCategoryDTO) *entity.BookEntity {
	id := uuid.UUID(b.ID.Bytes)
	var a []entity.BookAuthorEntity
	var c []entity.BookCategoryEntity

	for _, author := range authors {
		aID := uuid.UUID(author.ID.Bytes)
		a = append(a, entity.BookAuthorEntity{
			ID:   &aID,
			Name: author.Name,
		})
	}

	for _, category := range categories {
		cID := uuid.UUID(category.ID.Bytes)
		c = append(c, entity.BookCategoryEntity{
			ID:   &cID,
			Name: category.Name,
		})
	}

	return &entity.BookEntity{
		ID:            &id,
		Title:         b.Title,
		Subtitle:      b.Subtitle,
		Description:   b.Description,
		Authors:       a,
		Categories:    c,
		Year:          int16(b.PubYear),
		AverageRating: b.Rating,
		Embedding:     b.Embedding,
	}
}
