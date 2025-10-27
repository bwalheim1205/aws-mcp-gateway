package handlers

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type LambdaInvoker struct {
	client *lambda.Client
}

// NewLambdaInvoker loads AWS config and creates a Lambda client
func NewLambdaInvoker() (*LambdaInvoker, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	return &LambdaInvoker{
		client: lambda.NewFromConfig(cfg),
	}, nil
}

// InvokeLambda executes the Lambda function with the given name and arguments
func (li *LambdaInvoker) InvokeLambda(ctx context.Context, functionName string, args map[string]interface{}) (map[string]interface{}, error) {
	if li.client == nil {
		return nil, errors.New("lambda client not initialized")
	}

	payload, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	resp, err := li.client.Invoke(ctx, &lambda.InvokeInput{
		FunctionName:   aws.String(functionName),
		Payload:        payload,
		InvocationType: types.InvocationTypeRequestResponse,
	})
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Payload, &result); err != nil {
		return nil, err
	}

	return result, nil
}
