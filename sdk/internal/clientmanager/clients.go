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
)

type AWSClients struct {
	LambdaClient     *lambda.Client
	CloudWatchClient *cloudwatch.Client
	XRayClient       *xray.Client
	LogsClient       *cloudwatchlogs.Client
}

// NewAWSClients creates and returns AWS service clients required for the analyzer.
// 
// By default, it loads configuration using the AWS SDK's standard credential and
// configuration chain (environment variables, shared config files, EC2 instance roles, etc.).
//
// You can override the defaults by providing custom LoadOptions via the opts parameter,
// such as specifying a region, credentials provider, or profile.
//
// Returns an AWSClients struct containing initialized clients or an error if config loading fails.
func NewAWSClients(ctx context.Context, opts ...func(*config.LoadOptions) error) (*AWSClients, error) {
	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &AWSClients{
		LambdaClient:     lambda.NewFromConfig(cfg),
		CloudWatchClient: cloudwatch.NewFromConfig(cfg),
		XRayClient:       xray.NewFromConfig(cfg),
		LogsClient:       cloudwatchlogs.NewFromConfig(cfg),
	}, nil
}