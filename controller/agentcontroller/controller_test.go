package agentcontroller

import (
	"ai-agent/model/servicemodels"
	mockservicemodels "ai-agent/model/servicemodels/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAgentService := mockservicemodels.NewMockAgentService(t)
	controller := New(mockAgentService)

	t.Run("When request body is invalid, return bad request", func(t *testing.T) {
		router := gin.New()
		router.POST("/message", controller.SendMessage)

		req, _ := http.NewRequest(http.MethodPost, "/message", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "invalid request body", "code": "1"}`, w.Body.String())
	})

	t.Run("When user is not authenticated, return unauthorized", func(t *testing.T) {
		router := gin.New()
		router.POST("/message", controller.SendMessage)

		requestBody := servicemodels.RequestBody{
			Message: "Hello",
			UserID:  "",
		}
		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest(http.MethodPost, "/message", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"error": "user not authenticated"}`, w.Body.String())
	})

	t.Run("When service returns an error, return bad request", func(t *testing.T) {
		router := gin.New()
		router.POST("/message", controller.SendMessage)

		requestBody := servicemodels.RequestBody{
			Message: "Hello",
			UserID:  "test-user",
		}
		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest(http.MethodPost, "/message", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		expectedErr := errors.New("service error")
		mockAgentService.On("SendMessageWithHistory", "test-user", "Hello").Return("", expectedErr).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "service error"}`, w.Body.String())
	})

	t.Run("When successful, return created status and response", func(t *testing.T) {
		router := gin.New()
		router.POST("/message", controller.SendMessage)

		requestBody := servicemodels.RequestBody{
			Message: "Hello",
			UserID:  "test-user",
		}
		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest(http.MethodPost, "/message", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		serviceResponse := "Hi there!"
		mockAgentService.On("SendMessageWithHistory", "test-user", "Hello").Return(serviceResponse, nil).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, `"Hi there!"`, w.Body.String())
	})
}
