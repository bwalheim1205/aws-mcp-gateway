package tools

import (
	"context"
	"testing"

	"github.com/bwalheim1205/aws-mcp-gateway/internal/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func connectServerAndClient(t *testing.T) (*mcp.Server, *mcp.ClientSession) {
	t.Helper()

	server := mcp.NewServer(&mcp.Implementation{Name: "test-server", Version: "v0.1"}, nil)
	t1, t2 := mcp.NewInMemoryTransports()

	_, err := server.Connect(context.Background(), t1, nil)
	if err != nil {
		t.Fatalf("failed to connect server: %v", err)
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "v0.1"}, nil)
	sess, err := client.Connect(context.Background(), t2, nil)
	if err != nil {
		t.Fatalf("failed to connect client: %v", err)
	}

	return server, sess
}

func TestRegisterAndListTools(t *testing.T) {
	// Updated input schema
	inputSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"city": map[string]any{
				"type":        "string",
				"description": "City for which to get the weather",
			},
		},
		"required": []any{"city"},
	}

	cfg := &config.Config{
		Tools: []config.ToolDef{
			{
				Name:        "WeatherTool",
				Description: "Gets weather info",
				InputSchema: inputSchema,
			},
		},
	}

	server, clientSession := connectServerAndClient(t)

	if err := Register(server, cfg); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	resp, err := clientSession.ListTools(context.Background(), &mcp.ListToolsParams{Cursor: ""})
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}

	if len(resp.Tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(resp.Tools))
	}

	tool := resp.Tools[0]

	if tool.Name != "WeatherTool" {
		t.Errorf("expected tool name WeatherTool, got %s", tool.Name)
	}
	if tool.Description != "Gets weather info" {
		t.Errorf("description mismatch: %s", tool.Description)
	}

	// Validate schema exists
	if tool.InputSchema == nil {
		t.Fatalf("inputSchema is nil")
	}

	schemaMap, ok := tool.InputSchema.(map[string]any)
	if !ok {
		t.Fatalf("InputSchema is not a map[string]any: %#v", tool.InputSchema)
	}

	// Access the required field
	requiredField, ok := schemaMap["required"].([]any)
	if !ok {
		t.Fatalf("required field is not a list: %#v", schemaMap["required"])
	}

	if len(requiredField) != 1 || requiredField[0] != "city" {
		t.Fatalf("required list mismatch: %#v", requiredField)
	}
}

func TestRegisterEmptyConfig(t *testing.T) {
	cfg := &config.Config{Tools: []config.ToolDef{}}

	server, clientSession := connectServerAndClient(t)

	if err := Register(server, cfg); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	resp, err := clientSession.ListTools(context.Background(), &mcp.ListToolsParams{Cursor: ""})
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}

	if len(resp.Tools) != 0 {
		t.Errorf("expected no tools, got %d", len(resp.Tools))
	}
}
