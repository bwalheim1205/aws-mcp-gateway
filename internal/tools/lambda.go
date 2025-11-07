package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bwalheim1205/aws-mcp-gateway/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var lambdaClient *lambda.Client

func init() {
	cfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Sprintf("unable to load AWS SDK config: %v", err))
	}
	lambdaClient = lambda.NewFromConfig(cfg)
}

func LambdaHandler(td config.ToolDef) func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		payload, err := json.Marshal(req.Params)
		if err != nil {
			// Return the error directly
			return nil, fmt.Errorf("failed to marshal params: %w", err)
		}

		out, err := lambdaClient.Invoke(ctx, &lambda.InvokeInput{
			FunctionName:   aws.String(td.LambdaARN),
			Payload:        payload,
			InvocationType: types.InvocationTypeRequestResponse,
		})
		if err != nil {
			return nil, fmt.Errorf("lambda invoke failed: %w", err)
		}

		// Successful result
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: string(out.Payload),
				},
			},
		}, nil
	}
}
