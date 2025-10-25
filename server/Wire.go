//go:build wireinject
// +build wireinject

package server

import (
	"ai-agent/controller"
	"ai-agent/router"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func NewGinEngine() *gin.Engine {
	return gin.New()
}

var ProviderSet = wire.NewSet(
	NewGinEngine,
	controller.New,
	router.New,
	server.New,
)

func InitializeServer() *server.Server {
	panic(wire.Build(ProviderSet))

	return nil
}
