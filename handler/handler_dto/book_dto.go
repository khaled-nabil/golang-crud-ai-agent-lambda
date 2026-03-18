package handler_dto

import (
	"ai-agent/entity"
)

type (
	BookRequest struct {
		Title         string   `json:"title"`
		Subtitle      string   `json:"subtitle"`
		Authors       []string `json:"authors"`
		Categories    []string `json:"categories"`
		Thumbnail     string   `json:"thumbnail"`
		Description   string   `json:"description"`
		Year          int16    `json:"published_year"`
		AverageRating float32  `json:"average_rating"`
		RatingsCount  int      `json:"ratings_count"`
		PagesCount    int      `json:"num_pages"`
	}

	BookRecommendation struct {
		Title   string `json:"title"`
		Author  string `json:"author"`
		Genre   string `json:"genre"`
		Summary string `json:"summary"`
		Reason  string `json:"reason"`
	}

	RecommendResponseDTO struct {
		Description      string               `json:"description"`
		RecommendedBooks []BookRecommendation `json:"recommended_books"`
	}
)

func (r *BookRequest) ToBookEntity() *entity.BookEntity {
	var authors []entity.BookAuthorEntity
	var categories []entity.BookCategoryEntity

	for _, author := range r.Authors {
		authors = append(authors, entity.BookAuthorEntity{
			Name: author,
		})
	}

	for _, category := range r.Categories {
		categories = append(categories, entity.BookCategoryEntity{
			Name: category,
		})
	}

	return &entity.BookEntity{
		Title:         r.Title,
		Subtitle:      r.Subtitle,
		Authors:       authors,
		Categories:    categories,
		Description:   r.Description,
		Year:          r.Year,
		AverageRating: r.AverageRating,
		RatingCount:   r.RatingsCount,
		PageCount:     r.PagesCount,
	}
}
