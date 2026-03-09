package server

import (
	"github.com/gin-gonic/gin"

	"ai-agent/router"
)

type Server struct {
	engine *gin.Engine
	router *router.Router
}

func New(e *gin.Engine, r *router.Router) *Server {
	return &Server{e, r}
}

func (s *Server) Run(addr string) error {
	s.router.Route()
	return s.engine.Run(addr)
}
