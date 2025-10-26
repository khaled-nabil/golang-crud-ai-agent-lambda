package server

import (
	"context"

	"ai-agent/router"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	router *router.Router
}

func New(e *gin.Engine, r *router.Router) *Server {
	return &Server{e, r}
}

func (s *Server) Handle(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	e := ginadapter.New(s.engine)

	s.router.Route()

	return e.ProxyWithContext(ctx, *req)
}
