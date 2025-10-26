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

		s, err := server.InitializeServer()
		if err != nil {
			log.Fatalf("Failed to initialize server: %v", err)
		}

		return s.Handle(ctx, &req)
	})
}
