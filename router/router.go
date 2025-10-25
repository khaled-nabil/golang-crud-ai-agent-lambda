package router

import (
	"ai-agent/controller"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Gin  *gin.Engine
	ctrl *controller.Controller
}

func New(
	gin *gin.Engine,
	ctrl *controller.Controller,
) *Router {
	return &Router{gin, ctrl}
}

func (r *Router) Setup() {
	r.Gin.Group("/api/v1/").
		GET("health", r.ctrl.Health)
}
