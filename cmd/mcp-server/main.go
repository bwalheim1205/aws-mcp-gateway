package main

import (
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bwalheim1205/aws-mcp-gateway/internal/config"
	"github.com/bwalheim1205/aws-mcp-gateway/internal/tools"
)

func main() {
	cfg, err := config.LoadConfig("tools.yaml")
	if err != nil {
		log.Fatalf("failed to load tools: %v", err)
	}

	server := mcp.NewServer(&mcp.Implementation{Name: "LambdaMCPGateway", Version: "1.0.0"}, nil)

	tools.Register(server, cfg)

	// SSE endpoint – a simple implementation
	handler := mcp.NewSSEHandler(func(r *http.Request) *mcp.Server {
		return server
	}, nil)
	http.Handle("/mcp/sse", handler)

	// Serve http server
	log.Fatal(http.ListenAndServe(":8080", nil))

}
