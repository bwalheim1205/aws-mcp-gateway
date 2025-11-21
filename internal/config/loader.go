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
	Tools []ToolDef `yaml:"tools"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
