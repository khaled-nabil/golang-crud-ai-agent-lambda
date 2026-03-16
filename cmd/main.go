package main

import (
	"log"
	"os"

	"ai-agent/adapters/gemini"
	"ai-agent/adapters/secrets"
	"ai-agent/handler"
	"ai-agent/repositories/db"
	"ai-agent/router"
	"ai-agent/server"
	"ai-agent/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(os.Getenv("GIN_MODE"))
	engine := gin.New()

	appConfig, err := secrets.NewSecretsManager()
	if err != nil {
		log.Fatalf("Failed to load secrets: %v", err)
	}

	geminiClient, err := gemini.NewGeminiAdapter(appConfig)
	if err != nil {
		log.Fatalf("Failed to initialize Gemini: %v", err)
	}

	repository, err := db.NewPostgresRepo(appConfig)
	if err != nil {
		log.Fatalf("Failed to initialize DB: %v", err)
	}

	service := usecase.NewAIAgentUsecase(geminiClient, repository)
	errorHandler := usecase.NewErrorHandler()

	hh := handler.NewHealthHandler()
	ch := handler.NewChatHandler(service, errorHandler)

	r := router.New(engine, hh, ch)
	s := server.New(engine, r)

	if err := s.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
