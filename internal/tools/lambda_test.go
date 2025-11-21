package tools_test

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/bwalheim1205/aws-mcp-gateway/internal/config"
	"github.com/bwalheim1205/aws-mcp-gateway/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Mock client
type MockLambdaClient struct {
	Response *lambda.InvokeOutput
	Err      error
}

func (m *MockLambdaClient) Invoke(ctx context.Context, input *lambda.InvokeInput, optFns ...func(*lambda.Options)) (*lambda.InvokeOutput, error) {
	return m.Response, m.Err
}

func TestLambdaHandlerWithClient_Success(t *testing.T) {
	mockClient := &MockLambdaClient{
		Response: &lambda.InvokeOutput{
			Payload: []byte(`"hello world"`),
		},
	}

	td := config.ToolDef{
		LambdaARN: "arn:aws:lambda:us-east-1:123456789:function:test",
	}

	handler := tools.LambdaHandlerWithClient(td, mockClient)

	// Prepare raw arguments JSON
	args := map[string]interface{}{"foo": "bar"}
	rawArgs, err := json.Marshal(args)
	if err != nil {
		t.Fatalf("failed to marshal test args: %v", err)
	}

	req := &mcp.CallToolRequest{
		Params: &mcp.CallToolParamsRaw{
			Name:      td.Name, // set if needed
			Arguments: rawArgs,
		},
	}

	res, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Content) != 1 {
		t.Fatalf("expected 1 content element, got %d", len(res.Content))
	}

	textContent, ok := res.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected *mcp.TextContent, got %v", reflect.TypeOf(res.Content[0]))
	}

	if textContent.Text != `"hello world"` {
		t.Fatalf("unexpected text: %s", textContent.Text)
	}
}

func TestLambdaHandlerWithClient_InvokeError(t *testing.T) {
	mockClient := &MockLambdaClient{
		Err: errors.New("lambda failed"),
	}

	td := config.ToolDef{
		LambdaARN: "arn:aws:lambda:us-east-1:123456789:function:test",
	}

	handler := tools.LambdaHandlerWithClient(td, mockClient)

	args := map[string]interface{}{"foo": "bar"}
	rawArgs, err := json.Marshal(args)
	if err != nil {
		t.Fatalf("failed to marshal test args: %v", err)
	}

	req := &mcp.CallToolRequest{
		Params: &mcp.CallToolParamsRaw{
			Name:      td.Name,
			Arguments: rawArgs,
		},
	}

	res, err := handler(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if res != nil {
		t.Fatalf("expected nil result on error, got %v", res)
	}
}

func TestLambdaHandlerWithClient_MarshalError(t *testing.T) {
	mockClient := &MockLambdaClient{}

	td := config.ToolDef{
		LambdaARN: "arn:aws:lambda:us-east-1:123456789:function:test",
	}

	handler := tools.LambdaHandlerWithClient(td, mockClient)

	// invalid args that cannot be marshaled
	// but since you prepare rawArgs, you might instead pass invalid JSON intentionally:
	rawArgs := []byte{0xff, 0xfe, 0xfd} // invalid JSON

	req := &mcp.CallToolRequest{
		Params: &mcp.CallToolParamsRaw{
			Name:      td.Name,
			Arguments: rawArgs,
		},
	}

	res, err := handler(context.Background(), req)
	if err == nil {
		t.Fatal("expected marshal error, got nil")
	}
	if res != nil {
		t.Fatalf("expected nil result, got %v", res)
	}
}
