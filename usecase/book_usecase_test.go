package usecase

import (
	"ai-agent/entity"
	interfacemocks "ai-agent/interface/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBookUsecase_Insert(t *testing.T) {
	book := &entity.BookEntity{
		Title:       "Go Programming",
		Subtitle:    "A Modern Introduction",
		Description: "Learn Go",
		Authors:     []entity.BookAuthorEntity{{Name: "John"}},
		Categories:  []entity.BookCategoryEntity{{Name: "Technology"}},
	}
	expectedDocText := "title: Go Programming.\nsubtitle: A Modern Introduction.\ndescription: Learn Go.\nby: John, .\ncategory: Technology, ."

	t.Run("success", func(t *testing.T) {
		emMock := interfacemocks.NewMockEmbeddingModel(t)
		repoMock := interfacemocks.NewMockBookRepo(t)
		u := NewBookUsecase(emMock, repoMock)

		embedding := []float32{0.1, 0.2, 0.3}
		emMock.EXPECT().EmbedSearchDocument(expectedDocText).Return(embedding, nil)
		repoMock.EXPECT().InsertBookIfNotExists(book).Return(nil)

		err := u.Insert(book)

		require.NoError(t, err)
		require.Equal(t, embedding, book.Embedding)
	})

	t.Run("embed error", func(t *testing.T) {
		emMock := interfacemocks.NewMockEmbeddingModel(t)
		repoMock := interfacemocks.NewMockBookRepo(t)
		u := NewBookUsecase(emMock, repoMock)

		emMock.EXPECT().EmbedSearchDocument(expectedDocText).Return(nil, errors.New("embed failed"))

		err := u.Insert(book)

		require.Error(t, err)
		require.ErrorContains(t, err, "embed failed")
	})

	t.Run("insert error", func(t *testing.T) {
		emMock := interfacemocks.NewMockEmbeddingModel(t)
		repoMock := interfacemocks.NewMockBookRepo(t)
		u := NewBookUsecase(emMock, repoMock)

		embedding := []float32{0.1, 0.2, 0.3}
		emMock.EXPECT().EmbedSearchDocument(expectedDocText).Return(embedding, nil)
		repoMock.EXPECT().InsertBookIfNotExists(book).Return(errors.New("db error"))

		err := u.Insert(book)

		require.Error(t, err)
		require.ErrorContains(t, err, "db error")
	})
}

func TestBookUsecase_GetBookRecommendations(t *testing.T) {
	prompt := "science fiction books"

	t.Run("success", func(t *testing.T) {
		emMock := interfacemocks.NewMockEmbeddingModel(t)
		repoMock := interfacemocks.NewMockBookRepo(t)
		u := NewBookUsecase(emMock, repoMock)

		embedding := []float32{0.5, 0.6}
		books := []entity.BookEntity{
			{Title: "Dune"},
			{Title: "Foundation"},
		}
		emMock.EXPECT().EmbedSearchQuery(prompt).Return(embedding, nil)
		repoMock.EXPECT().SearchForRelevantBook(embedding).Return(books, nil)

		result, err := u.GetBookRecommendations(prompt)

		require.NoError(t, err)
		require.Len(t, result, 2)
		require.Equal(t, "Dune", result[0].Title)
	})

	t.Run("embed error", func(t *testing.T) {
		emMock := interfacemocks.NewMockEmbeddingModel(t)
		repoMock := interfacemocks.NewMockBookRepo(t)
		u := NewBookUsecase(emMock, repoMock)

		emMock.EXPECT().EmbedSearchQuery(prompt).Return(nil, errors.New("embed failed"))

		result, err := u.GetBookRecommendations(prompt)

		require.Error(t, err)
		require.Nil(t, result)
		require.ErrorContains(t, err, "embed failed")
	})

	t.Run("search error", func(t *testing.T) {
		emMock := interfacemocks.NewMockEmbeddingModel(t)
		repoMock := interfacemocks.NewMockBookRepo(t)
		u := NewBookUsecase(emMock, repoMock)

		embedding := []float32{0.5, 0.6}
		emMock.EXPECT().EmbedSearchQuery(prompt).Return(embedding, nil)
		repoMock.EXPECT().SearchForRelevantBook(embedding).Return(nil, errors.New("search failed"))

		result, err := u.GetBookRecommendations(prompt)

		require.Error(t, err)
		require.Nil(t, result)
		require.ErrorContains(t, err, "search failed")
	})
}
