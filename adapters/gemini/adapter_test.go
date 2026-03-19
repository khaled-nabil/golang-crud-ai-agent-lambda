package gemini

import (
	"ai-agent/entity"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/genai"
)

func TestGeminiAdapter_Chat_EmptyUserInput(t *testing.T) {
	adapter := &GeminiAdapter{}

	resp, err := adapter.Chat("", "sys", nil)

	require.Nil(t, resp)
	require.Error(t, err)
	require.ErrorContains(t, err, "user input or response is empty")
}

func TestTransformHistoryToGeminiContent(t *testing.T) {
	history := []entity.ChatHistoryEntity{
		{UserInput: "question-1", Response: "answer-1"},
		{UserInput: "question-2"},
		{Response: "answer-3"},
		{UserInput: "", Response: ""},
	}

	result := transformHistoryToGeminiContent(history)

	require.Len(t, result, 4)
	require.Equal(t, genai.RoleUser, result[0].Role)
	require.Equal(t, "question-1", result[0].Parts[0].Text)
	require.Equal(t, genai.RoleModel, result[1].Role)
	require.Equal(t, "answer-1", result[1].Parts[0].Text)
	require.Equal(t, genai.RoleUser, result[2].Role)
	require.Equal(t, "question-2", result[2].Parts[0].Text)
	require.Equal(t, genai.RoleModel, result[3].Role)
	require.Equal(t, "answer-3", result[3].Parts[0].Text)
}

func TestTransformBookToPrompt(t *testing.T) {
	bookID := uuid.New()
	book := &entity.BookEntity{
		ID:          &bookID,
		Title:       "The Pragmatic Programmer",
		Description: "A practical book about software craftsmanship",
		Authors: []entity.BookAuthorEntity{
			{Name: "Andrew Hunt"},
			{Name: "David Thomas"},
		},
		Categories: []entity.BookCategoryEntity{
			{Name: "Software"},
			{Name: "Engineering"},
		},
	}

	prompt := transformBookToPrompt(book)

	require.Contains(t, prompt, bookID.String())
	require.Contains(t, prompt, "Title: The Pragmatic Programmer")
	require.Contains(t, prompt, "Author: Andrew Hunt, David Thomas, ")
	require.Contains(t, prompt, "Genre: Software, Engineering, ")
	require.Contains(t, prompt, "Summary: A practical book about software craftsmanship")
}
