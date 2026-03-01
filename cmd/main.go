package main

import (
	"ai-agent/server"
	"log"
)

func main() {
	log.Print("Starting Server")

	s, err := server.InitializeServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	if err := s.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
