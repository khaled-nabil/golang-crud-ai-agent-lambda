package router

import (
	"ai-agent/controller/health"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Gin  *gin.Engine
	ctrl *health.Controller
}

func New(
	gin *gin.Engine,
	ctrl *health.Controller,
) *Router {
	return &Router{gin, ctrl}
}

func (r *Router) Route() {
	r.Gin.Group("/api/v1").
		GET("/health", r.ctrl.Health)
}
