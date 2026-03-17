package router

import (
	"ai-agent/handler"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Gin *gin.Engine
	hh  *handler.HealthHandler
	ch  *handler.ChatHandler
	bh  *handler.BookHandler
}

func New(
	gin *gin.Engine,
	hh *handler.HealthHandler,
	ch *handler.ChatHandler,
	bh *handler.BookHandler,
) *Router {
	return &Router{gin, hh, ch, bh}
}

func (r *Router) Route() {
	v1group := r.Gin.Group("/api/v1")

	v1group.
		Group("/health").
		GET("", r.hh.Health)

	v1group.
		Group("/chat").
		POST("", r.ch.ChatWithHistory)

	v1group.
		Group("/book").
		POST("", r.bh.CreateBook)
}
