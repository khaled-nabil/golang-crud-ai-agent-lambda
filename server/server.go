package server

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
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

func (s *Server) Handle(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	e := ginadapter.New(s.engine)

	s.router.Route()

	return e.ProxyWithContext(ctx, *req)
}
