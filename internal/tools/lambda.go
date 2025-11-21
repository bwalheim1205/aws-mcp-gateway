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

var globalLambdaClient *lambda.Client

func init() {
	cfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Sprintf("unable to load AWS SDK config: %v", err))
	}
	globalLambdaClient = lambda.NewFromConfig(cfg)
}

type LambdaClient interface {
	Invoke(ctx context.Context, params *lambda.InvokeInput, optFns ...func(*lambda.Options)) (*lambda.InvokeOutput, error)
}

func LambdaHandler(td config.ToolDef) func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return LambdaHandlerWithClient(td, nil)
}

func LambdaHandlerWithClient(td config.ToolDef, client LambdaClient) func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if client == nil {
		client = globalLambdaClient
	}

	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		payload, err := json.Marshal(req.Params.Arguments)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal params: %w", err)
		}

		out, err := client.Invoke(ctx, &lambda.InvokeInput{
			FunctionName:   aws.String(td.LambdaARN),
			Payload:        payload,
			InvocationType: types.InvocationTypeRequestResponse,
		})
		if err != nil {
			return nil, fmt.Errorf("lambda invoke failed: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: string(out.Payload),
				},
			},
		}, nil
	}
}
