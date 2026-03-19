package usecase

import (
	"ai-agent/entity"
	interfacemocks "ai-agent/interface/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAIAgentUsecase_SendMessageWithHistory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		agentMock := interfacemocks.NewMockAgentAdapter(t)
		chatRepoMock := interfacemocks.NewMockChatRepo(t)
		u := NewAIAgentUsecase(agentMock, chatRepoMock)

		userID := "user-1"
		message := "hello"
		systemPrompt := "sys"
		inputEmbedding := []float32{0.1, 0.2}
		history := []entity.ChatHistoryEntity{{UserInput: "old q", Response: "old a"}}
		chatResponse := &entity.ChatHistoryEntity{UserInput: message, Response: "new answer"}
		conversationEmbedding := []float32{0.9, 0.8}

		agentMock.EXPECT().EmbedMessage(message).Return(inputEmbedding, nil)
		chatRepoMock.EXPECT().GetUserSimilarDocuments(userID, inputEmbedding).Return(history, nil)
		agentMock.EXPECT().Chat(message, systemPrompt, history).Return(chatResponse, nil)
		agentMock.EXPECT().EmbedConverastion(chatResponse).Return(conversationEmbedding, nil)
		chatRepoMock.EXPECT().StoreConversation(userID, chatResponse, conversationEmbedding).Return(chatResponse, nil)

		resp, err := u.SendMessageWithHistory(userID, message, systemPrompt)

		require.NoError(t, err)
		require.Equal(t, "new answer", resp)
	})

	t.Run("embed message error", func(t *testing.T) {
		agentMock := interfacemocks.NewMockAgentAdapter(t)
		chatRepoMock := interfacemocks.NewMockChatRepo(t)
		u := NewAIAgentUsecase(agentMock, chatRepoMock)

		agentMock.EXPECT().EmbedMessage("hello").Return(nil, errors.New("embed failed"))

		resp, err := u.SendMessageWithHistory("user-1", "hello", "sys")

		require.Error(t, err)
		require.Empty(t, resp)
		require.ErrorContains(t, err, "failed to embed user input")
	})

	t.Run("get user similar documents error", func(t *testing.T) {
		agentMock := interfacemocks.NewMockAgentAdapter(t)
		chatRepoMock := interfacemocks.NewMockChatRepo(t)
		u := NewAIAgentUsecase(agentMock, chatRepoMock)

		inputEmbedding := []float32{0.1}
		agentMock.EXPECT().EmbedMessage("hello").Return(inputEmbedding, nil)
		chatRepoMock.EXPECT().GetUserSimilarDocuments("user-1", inputEmbedding).Return(nil, errors.New("db failed"))

		resp, err := u.SendMessageWithHistory("user-1", "hello", "sys")

		require.Error(t, err)
		require.Empty(t, resp)
		require.ErrorContains(t, err, "failed to retrieve user similar documents")
	})

	t.Run("chat error", func(t *testing.T) {
		agentMock := interfacemocks.NewMockAgentAdapter(t)
		chatRepoMock := interfacemocks.NewMockChatRepo(t)
		u := NewAIAgentUsecase(agentMock, chatRepoMock)

		inputEmbedding := []float32{0.1}
		history := []entity.ChatHistoryEntity{{UserInput: "old q", Response: "old a"}}
		agentMock.EXPECT().EmbedMessage("hello").Return(inputEmbedding, nil)
		chatRepoMock.EXPECT().GetUserSimilarDocuments("user-1", inputEmbedding).Return(history, nil)
		agentMock.EXPECT().Chat("hello", "sys", history).Return(nil, errors.New("chat failed"))

		resp, err := u.SendMessageWithHistory("user-1", "hello", "sys")

		require.Error(t, err)
		require.Empty(t, resp)
		require.ErrorContains(t, err, "failed to send message to user")
	})

	t.Run("embed conversation error", func(t *testing.T) {
		agentMock := interfacemocks.NewMockAgentAdapter(t)
		chatRepoMock := interfacemocks.NewMockChatRepo(t)
		u := NewAIAgentUsecase(agentMock, chatRepoMock)

		inputEmbedding := []float32{0.1}
		history := []entity.ChatHistoryEntity{{UserInput: "old q", Response: "old a"}}
		chatResponse := &entity.ChatHistoryEntity{UserInput: "hello", Response: "answer"}
		agentMock.EXPECT().EmbedMessage("hello").Return(inputEmbedding, nil)
		chatRepoMock.EXPECT().GetUserSimilarDocuments("user-1", inputEmbedding).Return(history, nil)
		agentMock.EXPECT().Chat("hello", "sys", history).Return(chatResponse, nil)
		agentMock.EXPECT().EmbedConverastion(chatResponse).Return(nil, errors.New("embed conversation failed"))

		resp, err := u.SendMessageWithHistory("user-1", "hello", "sys")

		require.Error(t, err)
		require.Empty(t, resp)
		require.ErrorContains(t, err, "failed to embed conversation")
	})

	t.Run("store conversation error", func(t *testing.T) {
		agentMock := interfacemocks.NewMockAgentAdapter(t)
		chatRepoMock := interfacemocks.NewMockChatRepo(t)
		u := NewAIAgentUsecase(agentMock, chatRepoMock)

		inputEmbedding := []float32{0.1}
		history := []entity.ChatHistoryEntity{{UserInput: "old q", Response: "old a"}}
		chatResponse := &entity.ChatHistoryEntity{UserInput: "hello", Response: "answer"}
		conversationEmbedding := []float32{0.9}
		agentMock.EXPECT().EmbedMessage("hello").Return(inputEmbedding, nil)
		chatRepoMock.EXPECT().GetUserSimilarDocuments("user-1", inputEmbedding).Return(history, nil)
		agentMock.EXPECT().Chat("hello", "sys", history).Return(chatResponse, nil)
		agentMock.EXPECT().EmbedConverastion(chatResponse).Return(conversationEmbedding, nil)
		repoErr := errors.New("store failed")
		chatRepoMock.EXPECT().StoreConversation("user-1", chatResponse, conversationEmbedding).Return(nil, repoErr)

		resp, err := u.SendMessageWithHistory("user-1", "hello", "sys")

		require.Error(t, err)
		require.Empty(t, resp)
		require.ErrorIs(t, err, repoErr)
	})
}
