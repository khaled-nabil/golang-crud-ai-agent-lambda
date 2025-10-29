//go:build wireinject
// +build wireinject

package server

import (
	"ai-agent/controller/agentcontroller"
	"ai-agent/controller/healthcontroller"
	"ai-agent/model/geminimodel"
	"ai-agent/pkg/geminipkg"
	"ai-agent/pkg/secretspkg"
	"ai-agent/router"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func NewGinEngine() *gin.Engine {
	gin.SetMode(os.Getenv("GIN_MODE"))

	return gin.New()
}

var ProviderSet = wire.NewSet(
	NewGinEngine,
	New,
	healthcontroller.New,
	wire.Bind(new(geminimodel.Gemini), new(*geminipkg.Gemini)),
	agentcontroller.New,
	router.New,
	geminipkg.New,
	secretspkg.New,
)

func InitializeServer() (*Server, error) {
	wire.Build(ProviderSet)

	return nil, nil
}
