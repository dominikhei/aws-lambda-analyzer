package fetcherinterfaces

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

// This interface matches logsinsightsfetcher.Fetcher for tetsing the internal functions
type LogsInsightsFetcher interface {
	RunQuery(ctx context.Context, fq sdktypes.FunctionQuery, queryString string) ([]map[string]string, error)
}

// This interface matches cloudwatchfetcher.Fetcher for tetsing the internal functions
type CloudWatchFetcher interface {
	FetchMetric(ctx context.Context, query sdktypes.FunctionQuery, metricName string, stat string) ([]types.MetricDataResult, error)
}

// This interface matches lambda.Client for tetsing the internal functions
type LambdaClient interface {
	GetFunction(ctx context.Context, params *lambda.GetFunctionInput, optFns ...func(*lambda.Options)) (*lambda.GetFunctionOutput, error)
}
