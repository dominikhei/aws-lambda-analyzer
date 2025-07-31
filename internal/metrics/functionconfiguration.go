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

package metrics

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	sdkinterfaces "github.com/dominikhei/serverless-statistics/internal/interfaces"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

// GetFunctionConfiguration gets configurations of an AWS Lambda function with a sprcific qualifier.
func GetFunctionConfiguration(
	ctx context.Context,
	lambdaClient sdkinterfaces.LambdaClient,
	query sdktypes.FunctionQuery,
) (*sdktypes.BaseStatisticsReturn, error) {

	funcConfig, err := lambdaClient.GetFunction(ctx, &lambda.GetFunctionInput{
		FunctionName: aws.String(query.FunctionName),
		Qualifier:    aws.String(query.Qualifier),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get function configuration: %w", err)
	}

	envVars := make(map[string]string)
	if funcConfig.Configuration.Environment != nil && funcConfig.Configuration.Environment.Variables != nil {
		envVars = funcConfig.Configuration.Environment.Variables
	}
	return &sdktypes.BaseStatisticsReturn{
		FunctionARN:          aws.ToString(funcConfig.Configuration.FunctionArn),
		FunctionName:         aws.ToString(funcConfig.Configuration.FunctionName),
		Qualifier:            aws.ToString(funcConfig.Configuration.Version),
		MemorySizeMB:         funcConfig.Configuration.MemorySize,
		TimeoutSeconds:       funcConfig.Configuration.Timeout,
		Runtime:              string(funcConfig.Configuration.Runtime),
		LastModified:         aws.ToString(funcConfig.Configuration.LastModified),
		EnvironmentVariables: envVars,
	}, nil
}
