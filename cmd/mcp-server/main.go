package main

import (
	"log"

	"github.com/bwalheim1205/aws-mcp-gateway/internal/handlers"
	"github.com/bwalheim1205/aws-mcp-gateway/internal/server"
)

func main() {
	// Initialize Lambda invoker
	lambdaInvoker, err := handlers.NewLambdaInvoker()
	if err != nil {
		log.Fatal("Failed to initialize Lambda invoker:", err)
	}

	// Start HTTP server
	port := "8080"
	log.Printf("Starting MCP server on port %s...\n", port)
	if err := server.StartServer(port, lambdaInvoker); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
