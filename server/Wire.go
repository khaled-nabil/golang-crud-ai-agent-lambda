//go:build wireinject
// +build wireinject

package server

import (
	"ai-agent/adapters/gemini"
	"ai-agent/adapters/ollama"
	"ai-agent/adapters/secrets"
	"ai-agent/controller/agentcontroller"
	"ai-agent/controller/healthcontroller"
	"ai-agent/model/datamodels"
	"ai-agent/model/errormodels"
	"ai-agent/model/servicemodels"
	"ai-agent/repositories/db"
	"ai-agent/router"
	"ai-agent/service/aiagent"
	"ai-agent/usecase/errors"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func NewGinEngine() *gin.Engine {
	gin.SetMode(os.Getenv("GIN_MODE"))

	return gin.New()
}

var ProviderSet = wire.NewSet(
	wire.Bind(new(datamodels.Gemini), new(*gemini.Gemini)),
	wire.Bind(new(datamodels.EmbeddingModel), new(*ollama.Ollama)),
	wire.Bind(new(servicemodels.AgentService), new(*aiagent.Service)),
	wire.Bind(new(servicemodels.Persistence), new(*db.Repository)),
	wire.Bind(new(errormodels.Errors), new(*errors.ErrorHandler)),
	NewGinEngine,
	New,
	healthcontroller.New,
	agentcontroller.New,
	router.New,
	gemini.New,
	ollama.New,
	secrets.New,
	aiagent.New,
	db.New,
	errors.New,
)

func InitializeServer() (*Server, error) {
	wire.Build(ProviderSet)

	return nil, nil
}
