package tools

import (
	"github.com/bwalheim1205/aws-mcp-gateway/internal/config"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Register(server *mcp.Server, cfg *config.Config) error {
	for _, td := range cfg.Tools {

		tool := &mcp.Tool{
			Name:        td.Name,
			Description: td.Description,
			InputSchema: td.InputSchema,
		}

		// Register Lambda Tool
		server.AddTool(tool, LambdaHandler(td))
	}

	return nil
}
