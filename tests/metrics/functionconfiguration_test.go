// Copyright 2025 dominikhei
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tests

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/dominikhei/serverless-statistics/internal/metrics"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
	"github.com/stretchr/testify/require"
)

func TestGetFunctionConfiguration(t *testing.T) {
	mockLambdaClient := &mockLambdaClient{
		GetFunctionFunc: func(ctx context.Context, params *lambda.GetFunctionInput, optFns ...func(*lambda.Options)) (*lambda.GetFunctionOutput, error) {
			return &lambda.GetFunctionOutput{
				Configuration: &types.FunctionConfiguration{
					FunctionName: aws.String("my-lambda-fn"),
					FunctionArn:  aws.String("arn:aws:lambda:us-east-1:123456789012:function:my-lambda-fn"),
					Version:      aws.String("1"),
					MemorySize:   aws.Int32(512),
					Timeout:      aws.Int32(15),
					Runtime:      types.RuntimeGo1x,
					LastModified: aws.String("2023-01-01T00:00:00.000+0000"),
					Environment: &types.EnvironmentResponse{
						Variables: map[string]string{
							"ENV": "prod",
						},
					},
				},
			}, nil
		},
	}

	query := sdktypes.FunctionQuery{
		FunctionName: "my-lambda-fn",
		Qualifier:    "1",
	}

	result, err := metrics.GetFunctionConfiguration(context.Background(), mockLambdaClient, query)
	require.NoError(t, err)

	require.Equal(t, "my-lambda-fn", result.FunctionName)
	require.Equal(t, "arn:aws:lambda:us-east-1:123456789012:function:my-lambda-fn", result.FunctionARN)
	require.Equal(t, int32(512), *result.MemorySizeMB)
	require.Equal(t, int32(15), *result.TimeoutSeconds)
	require.Equal(t, "go1.x", result.Runtime)
	require.Equal(t, "2023-01-01T00:00:00.000+0000", result.LastModified)
	require.Equal(t, map[string]string{"ENV": "prod"}, result.EnvironmentVariables)
}

func TestGetFunctionConfiguration_NoEnvVars(t *testing.T) {
	mockLambdaClient := &mockLambdaClient{
		GetFunctionFunc: func(ctx context.Context, params *lambda.GetFunctionInput, optFns ...func(*lambda.Options)) (*lambda.GetFunctionOutput, error) {
			return &lambda.GetFunctionOutput{
				Configuration: &types.FunctionConfiguration{
					FunctionName: aws.String("my-lambda-fn"),
					FunctionArn:  aws.String("arn:aws:lambda:us-east-1:123456789012:function:my-lambda-fn"),
					Version:      aws.String("1"),
					MemorySize:   aws.Int32(512),
					Timeout:      aws.Int32(15),
					Runtime:      types.RuntimeGo1x,
					LastModified: aws.String("2023-01-01T00:00:00.000+0000"),
					Environment: &types.EnvironmentResponse{
						Variables: map[string]string{},
					},
				},
			}, nil
		},
	}

	query := sdktypes.FunctionQuery{
		FunctionName: "test-fn",
		Qualifier:    "1",
	}

	result, err := metrics.GetFunctionConfiguration(context.Background(), mockLambdaClient, query)
	require.NoError(t, err)
	require.Empty(t, result.EnvironmentVariables)
}

func TestGetFunctionConfiguration_MissingMemoryAndTimeout(t *testing.T) {
	mockLambdaClient := &mockLambdaClient{
		GetFunctionFunc: func(ctx context.Context, params *lambda.GetFunctionInput, optFns ...func(*lambda.Options)) (*lambda.GetFunctionOutput, error) {
			return &lambda.GetFunctionOutput{
				Configuration: &types.FunctionConfiguration{
					FunctionName: aws.String("my-lambda-fn"),
					FunctionArn:  aws.String("arn:aws:lambda:us-east-1:123456789012:function:my-lambda-fn"),
					Version:      aws.String("1"),
					// MemorySize and Timeout are omitted
					Runtime:      types.RuntimeGo1x,
					LastModified: aws.String("2023-01-01T00:00:00.000+0000"),
					Environment: &types.EnvironmentResponse{
						Variables: map[string]string{"ENV": "prod"},
					},
				},
			}, nil
		},
	}

	query := sdktypes.FunctionQuery{
		FunctionName: "my-lambda-fn",
		Qualifier:    "1",
	}

	result, err := metrics.GetFunctionConfiguration(context.Background(), mockLambdaClient, query)
	require.NoError(t, err)

	require.Equal(t, "my-lambda-fn", result.FunctionName)
	require.Equal(t, "arn:aws:lambda:us-east-1:123456789012:function:my-lambda-fn", result.FunctionARN)
	require.Equal(t, "1", result.Qualifier)
	require.Nil(t, result.MemorySizeMB)
	require.Nil(t, result.TimeoutSeconds)
	require.Equal(t, "go1.x", result.Runtime)
	require.Equal(t, "2023-01-01T00:00:00.000+0000", result.LastModified)
	require.Equal(t, map[string]string{"ENV": "prod"}, result.EnvironmentVariables)
}
