//go:build wireinject
// +build wireinject

package server

import (
	"ai-agent/controller/agentcontroller"
	"ai-agent/controller/healthcontroller"
	"ai-agent/model/datamodels"
	"ai-agent/model/servicemodels"
	"ai-agent/adapters/gemini"
	"ai-agent/adapters/secrets"
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
	wire.Bind(new(servicemodels.AgentService), new(*aiagent.Service)),
	wire.Bind(new(servicemodels.Persistence), new(*db.Repository)),
	NewGinEngine,
	New,
	errors.New,
	healthcontroller.New,
	agentcontroller.New,
	router.New,
	gemini.New,
	secrets.New,
	aiagent.New,
	db.New,
)

func InitializeServer() (*Server, error) {
	wire.Build(ProviderSet)

	return nil, nil
}
