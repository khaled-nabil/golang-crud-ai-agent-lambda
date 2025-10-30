package router

import (
	"ai-agent/controller/agentcontroller"
	"ai-agent/controller/healthcontroller"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Gin   *gin.Engine
	hctrl *healthcontroller.Controller
	actrl *agentcontroller.Controller
}

func New(
	gin *gin.Engine,
	hctrl *healthcontroller.Controller,
	actrl *agentcontroller.Controller,
) *Router {
	return &Router{gin, hctrl, actrl}
}

func (r *Router) Route() {
	r.Gin.Group("/api/v1").
		GET("/health", r.hctrl.Health).
		POST("/send", r.actrl.SendMessage)
}
