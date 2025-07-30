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

// Package clientmanager provides a centralized clientmanager which handles the
// various aws clients required for this sdk.
package clientmanager

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/xray"
	"github.com/dominikhei/serverless-statistics/internal/utils"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

// NewAWSClients will be called by serverlessstatistics.New() internally to set up connections.
func NewAWSClients(ctx context.Context, opts sdktypes.ConfigOptions) (*sdktypes.AWSClients, error) {
	loadOpts, err := utils.ToLoadOptions(opts)
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadDefaultConfig(ctx, loadOpts...)
	if err != nil {
		return nil, err
	}

	return &sdktypes.AWSClients{
		LambdaClient:     lambda.NewFromConfig(cfg),
		CloudWatchClient: cloudwatch.NewFromConfig(cfg),
		XRayClient:       xray.NewFromConfig(cfg),
		LogsClient:       cloudwatchlogs.NewFromConfig(cfg),
	}, nil
}
