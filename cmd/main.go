package main

import (
	"ai-agent/server"
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		log.Print("Starting Server")

		s := server.InitializeServer()

		return s.Handle(ctx, &req)
	})
}
