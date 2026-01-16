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

type ServerConfig struct {
	Name     string `yaml:"name"`
	Version  string `yaml:"version"`
	Mode     string `yaml:"mode"`
	Endpoint string `yaml:"endpoint"`
	Port     int    `yaml:"port"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
	Tools  []ToolDef    `yaml:"tools"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Server: ServerConfig{
			Name:     "LambdaMCPGateway",
			Version:  "v1.0.0",
			Mode:     "stream",
			Endpoint: "/mcp/",
			Port:     8080,
		},
	}

	// Parse YAML into cfg
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
