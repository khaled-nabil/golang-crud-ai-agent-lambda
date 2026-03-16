package router

import (
	"ai-agent/handler"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Gin *gin.Engine
	hh  *handler.HealthHandler
	ch  *handler.ChatHandler
}

func New(
	gin *gin.Engine,
	hctrl *handler.HealthHandler,
	actrl *handler.ChatHandler,
) *Router {
	return &Router{gin, hctrl, actrl}
}

func (r *Router) Route() {
	v1group := r.Gin.Group("/api/v1")

	v1group.
		Group("/health").
		GET("", r.hh.Health)

	v1group.
		Group("/chat").
		POST("", r.ch.ChatWithHistory)
}
