package repo_dto

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type (
	BookDTO struct {
		ID          *pgtype.UUID `json:"id,omitempty"`
		Title       string       `json:"title"`
		Subtitle    string       `json:"subtitle,omitempty"`
		Description string       `json:"description"`
		Embedding   float64      `json:"embedding"`
		Thumb       string       `json:"thumbnail"`
		PubYear     int8         `json:"published_year"`
		Rating      float32      `json:"average_rating"`
		RatingCount int          `json:"rating_count"`
		PageCount   int16        `json:"num_pages"`
		CreatedAt   *time.Time   `json:"created_at,omitempty"`
	}
)
