package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ToolDef struct {
	Name        string         `yaml:"name"`
	LambdaARN   string         `yaml:"lambdaArn"`
	Description string         `yaml:"description"`
	InputSchema map[string]any `yaml:"inputSchema"`
}

type Config struct {
	Name     string    `yaml:"name"`
	Version  string    `yaml:"version"`
	Endpoint string    `yaml:"endpoint"`
	Port     int       `yaml:"port"`
	Tools    []ToolDef `yaml:"tools"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Default Values
	cfg := &Config{
		Name:     "LambdaMCPGateway",
		Version:  "v1.0.0",
		Endpoint: "/mcp/sse",
		Port:     8080,
	}

	// Parse YAML into cfg
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
