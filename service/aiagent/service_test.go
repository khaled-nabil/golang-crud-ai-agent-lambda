package aiagent

import (
	"ai-agent/model/datamodels"
	mockdatamodels "ai-agent/model/datamodels/mocks"
	mockservicemodels "ai-agent/model/servicemodels/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMessageWithHistory(t *testing.T) {
	mockAgent := mockdatamodels.NewMockGemini(t)
	mockDb := mockservicemodels.NewMockAgentRepo(t)

	s := New(mockAgent, mockDb)

	userID := "test-user"
	message := "Hello"

	t.Run("When new user, create chat with no history and store conversation", func(t *testing.T) {
		var history []datamodels.HistoryContext
		response := &datamodels.HistoryContext{
			UserInput: message,
			Response:  "Hi there!",
		}
		mockDb.On("GetUserHistory", userID).Return(history, nil).Once()
		mockAgent.On("Chat", message, history).Return(response, nil).Once()
		mockDb.On("StoreConversation", userID, response).Return(nil).Once()

		resp, err := s.SendMessageWithHistory(userID, message)
		assert.NoError(t, err)
		assert.Equal(t, response.Response, resp)
	})

	t.Run("When existing user, create chat with history and store conversation", func(t *testing.T) {
		history := []datamodels.HistoryContext{
			{
				UserInput: "previous message",
				Response:  "previous response",
			},
		}
		response := &datamodels.HistoryContext{
			UserInput: message,
			Response:  "Hi there!",
		}
		mockDb.On("GetUserHistory", userID).Return(history, nil).Once()
		mockAgent.On("Chat", message, history).Return(response, nil).Once()
		mockDb.On("StoreConversation", userID, response).Return(nil).Once()

		resp, err := s.SendMessageWithHistory(userID, message)
		assert.NoError(t, err)
		assert.Equal(t, response.Response, resp)
	})

	t.Run("When getting user history fails, return error response", func(t *testing.T) {
		expectedErr := errors.New("db error")
		mockDb.On("GetUserHistory", userID).Return(nil, expectedErr).Once()

		_, err := s.SendMessageWithHistory(userID, message)
		assert.Error(t, err)
		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("When AI agent fails, return error response", func(t *testing.T) {
		var history []datamodels.HistoryContext
		expectedErr := errors.New("chat error")
		mockDb.On("GetUserHistory", userID).Return(history, nil).Once()
		mockAgent.On("Chat", message, history).Return(nil, expectedErr).Once()

		_, err := s.SendMessageWithHistory(userID, message)
		assert.Error(t, err)
		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("When storing conversation fails, return error response", func(t *testing.T) {
		var history []datamodels.HistoryContext
		response := &datamodels.HistoryContext{
			UserInput: message,
			Response:  "Hi there!",
		}
		expectedErr := errors.New("store error")
		mockDb.On("GetUserHistory", userID).Return(history, nil).Once()
		mockAgent.On("Chat", message, history).Return(response, nil).Once()
		mockDb.On("StoreConversation", userID, response).Return(expectedErr).Once()

		_, err := s.SendMessageWithHistory(userID, message)
		assert.Error(t, err)
		assert.ErrorIs(t, err, expectedErr)
	})
}
