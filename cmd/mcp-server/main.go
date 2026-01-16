package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/pflag"

	"github.com/bwalheim1205/aws-mcp-gateway/internal/config"
	"github.com/bwalheim1205/aws-mcp-gateway/internal/tools"
)

// Build information (populated at build time)
var version = "dev"
var commit = "none"
var buildTime = "unknown"

// Configuration values
var configFile string
var port int
var mode string
var endpoint string
var versionFlag bool

func argparse() {
	// Configuration File
	pflag.StringVarP(&configFile, "files", "f", "tools.yaml", "Path to the tools configuration file")

	// Mode
	pflag.StringVarP(&mode, "mode", "m", "", "Port for mcp server to use")

	// Port
	pflag.IntVarP(&port, "port", "p", 0, "Port for mcp server to use")

	// Endpoint
	pflag.StringVarP(&endpoint, "endpoint", "e", "", "Port for mcp server to use")

	// Version
	pflag.BoolVar(&versionFlag, "version", false, "Version information")

	//Parse
	pflag.Parse()
}

func main() {
	argparse()

	if versionFlag {
		fmt.Printf("Version: %s\nCommit: %s\nBuild Time: %s\n", version, commit, buildTime)
		return
	}

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("failed to load configuration file: %v", err)
	}

	// Load config if not specified on cli
	if port == 0 {
		port = cfg.Server.Port
	}
	if endpoint == "" {
		endpoint = cfg.Server.Endpoint
	}
	if mode == "" {
		mode = cfg.Server.Mode
	}

	server := mcp.NewServer(&mcp.Implementation{Name: cfg.Server.Name, Version: cfg.Server.Version}, nil)

	tools.Register(server, cfg)

	// SSE endpoint – a simple implementation
	switch mode {
	case "sse":
		http.Handle(endpoint, mcp.NewSSEHandler(func(r *http.Request) *mcp.Server {
			return server
		}, nil))
	case "stream":
		http.Handle(endpoint, mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
			return server
		}, nil))
	default:
		log.Fatalf("Unknown mode (%s)! Please update server.mode in config. Valid modes: sse, stream", mode)
	}

	log.Printf("Listening at http://localhost:%d%s", port, endpoint)

	// Serve http server
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))

}
