package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for test YAML files
	tmpDir := t.TempDir()

	validYAML := `
name: "MyGateway"
version: "1.2.3"
endpoint: "http://localhost"
port: 8080

tools:
- name: "ExampleTool"
  lambdaArn: "arn:aws:lambda:us-east-1:123456789012:function:Example"
  description: "This is an example tool"
  inputSchema:
    param1: string
    param2: int
  `

	// Write valid YAML to file
	validFile := filepath.Join(tmpDir, "valid.yaml")
	if err := os.WriteFile(validFile, []byte(validYAML), 0644); err != nil {
		t.Fatalf("failed to write valid YAML: %v", err)
	}

	t.Run("Load valid config", func(t *testing.T) {
		cfg, err := LoadConfig(validFile)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// ---- NEW TESTS FOR OPTIONAL FIELDS ----
		if cfg.Name != "MyGateway" {
			t.Errorf("expected name 'MyGateway', got %s", cfg.Name)
		}
		if cfg.Version != "1.2.3" {
			t.Errorf("expected version '1.2.3', got %s", cfg.Version)
		}
		if cfg.Endpoint != "http://localhost" {
			t.Errorf("expected endpoint 'http://localhost', got %s", cfg.Endpoint)
		}
		if cfg.Port != 8080 {
			t.Errorf("expected port 8080, got %d", cfg.Port)
		}

		// ---- EXISTING TOOL TESTS ----
		if len(cfg.Tools) != 1 {
			t.Fatalf("expected 1 tool, got %d", len(cfg.Tools))
		}
		tool := cfg.Tools[0]
		if tool.Name != "ExampleTool" {
			t.Errorf("expected tool name ExampleTool, got %s", tool.Name)
		}
		if tool.LambdaARN != "arn:aws:lambda:us-east-1:123456789012:function:Example" {
			t.Errorf("unexpected LambdaARN: %s", tool.LambdaARN)
		}
		if tool.Description != "This is an example tool" {
			t.Errorf("unexpected description: %s", tool.Description)
		}
		if len(tool.InputSchema) != 2 {
			t.Errorf("expected 2 schema entries, got %d", len(tool.InputSchema))
		}
	})

	t.Run("Load missing file", func(t *testing.T) {
		_, err := LoadConfig(filepath.Join(tmpDir, "missing.yaml"))
		if err == nil {
			t.Fatal("expected error for missing file, got nil")
		}
	})

	// ---- NEW TEST: Missing optional fields are allowed ----
	optionalMissingYAML := `
tools:
- name: "ExampleTool"
  lambdaArn: "arn:aws:lambda:us-east-1:123456789012:function:Example"
  description: "This is an example tool"
  inputSchema:
    param1: string
    param2: int
  `
	missingOptFile := filepath.Join(tmpDir, "missing_optional.yaml")
	if err := os.WriteFile(missingOptFile, []byte(optionalMissingYAML), 0644); err != nil {
		t.Fatalf("failed to write optional-missing YAML: %v", err)
	}

	t.Run("Load config with missing optional fields", func(t *testing.T) {
		cfg, err := LoadConfig(missingOptFile)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Optional fields should be zero-value
		if cfg.Name != "LambdaMCPGateway" {
			t.Errorf("expected default name 'LambdaMCPGateway', got %s", cfg.Name)
		}
		if cfg.Version != "v1.0.0" {
			t.Errorf("expected default version 'v1.0.0', got %s", cfg.Version)
		}
		if cfg.Endpoint != "/mcp/sse" {
			t.Errorf("expected default endpoint '/mcp/sse', got %s", cfg.Endpoint)
		}
		if cfg.Port != 8080 {
			t.Errorf("expected default port 8080, got %d", cfg.Port)
		}

		if len(cfg.Tools) != 1 {
			t.Fatalf("expected 1 tool, got %d", len(cfg.Tools))
		}
	})
}
