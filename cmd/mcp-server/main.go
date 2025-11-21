package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bwalheim1205/aws-mcp-gateway/internal/config"
	"github.com/bwalheim1205/aws-mcp-gateway/internal/tools"
)

func main() {
	configFile := flag.String("f", "tools.yaml", "Path to the tools configuration file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configFile)
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
