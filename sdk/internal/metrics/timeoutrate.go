package metrics

import (
    "context"
    "errors"
    "fmt"
    "strconv"

    cloudwatchfetcher "github.com/dominikhei/aws-lambda-analyzer/sdk/internal/cloudwatch"
    logsinsightsfetcher "github.com/dominikhei/aws-lambda-analyzer/sdk/internal/logsinsights"
    sdktypes "github.com/dominikhei/aws-lambda-analyzer/sdk/types"
	"github.com/dominikhei/aws-lambda-analyzer/sdk/internal/queries"
)

func GetTimeoutRate(
    ctx context.Context,
    cwFetcher *cloudwatchfetcher.Fetcher,
    logsFetcher *logsinsightsfetcher.Fetcher,
    query sdktypes.FunctionQuery,
    period int32,
) (*sdktypes.TimeoutRateReturn, error) {
    results, err := logsFetcher.RunQuery(ctx, query, queries.LambdaTimeoutQuery)
	fmt.Print((results))
    if err != nil {
        return nil, fmt.Errorf("run logs insights query: %w", err)
    }

    var timeoutCount int
    if len(results) > 0 {
        val := results[0]["timeoutCount"]
        if val != "" {
            timeoutCount, err = strconv.Atoi(val)
            if err != nil {
                return nil, fmt.Errorf("parse timeoutCount from logs: %w", err)
            }
        }
    }
    invocationsResults, err := cwFetcher.FetchMetric(ctx, query, "Invocations", "Sum", period)
    if err != nil {
        return nil, fmt.Errorf("fetch invocations metric: %w", err)
    }

    invocationsSum, err := sumMetricValues(invocationsResults)
    if err != nil {
        return nil, fmt.Errorf("parse invocations metric data: %w", err)
    }

    if invocationsSum == 0 {
        return nil, errors.New("total invocations is zero, cannot calculate timeout rate")
    }

    timeoutRate := float64(timeoutCount) / invocationsSum

    return &sdktypes.TimeoutRateReturn{
        TimeoutRate:  timeoutRate,
        FunctionName: query.FunctionName,
        Qualifier:    query.Qualifier,
        StartTime:    query.StartTime,
        EndTime:      query.EndTime,
    }, nil
}