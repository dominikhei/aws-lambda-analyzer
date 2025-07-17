package metrics

import (
	"github.com/dominikhei/aws-lambda-analyzer/sdk/internal/cloudwatch"
    "context"
    "errors"
    "fmt"
    "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
    sdktypes "github.com/dominikhei/aws-lambda-analyzer/sdk/types"
)

func GetThrottleRate(ctx context.Context, f *cloudwatchfetcher.Fetcher, query sdktypes.FunctionQuery, period int32) (*sdktypes.ThrottleRateReturn, error) {
    throttlesResults, err := f.FetchMetric(ctx, query, "Throttles", "Sum", period)
    if err != nil {
        return nil, fmt.Errorf("fetch throttles metric: %w", err)
    }
    invocationsResults, err := f.FetchMetric(ctx, query, "Invocations", "Sum", period)
    if err != nil {
        return nil, fmt.Errorf("fetch invocations metric: %w", err)
    }

    throttlesSum, err := sumMetricValues(throttlesResults)
    if err != nil {
        return nil, fmt.Errorf("parse throttles metric data: %w", err)
    }
    invocationsSum, err := sumMetricValues(invocationsResults)
    if err != nil {
        return nil, fmt.Errorf("parse invocations metric data: %w", err)
    }
    if invocationsSum == 0 {
        return nil, errors.New("total invocations is zero, cannot calculate throttle rate")
    }
    throttleRate := throttlesSum / invocationsSum
	result := &sdktypes.ThrottleRateReturn{
        ThrottleRate: throttleRate,
        FunctionName: query.FunctionName,
        Qualifier:    query.Qualifier,
        StartTime:    query.StartTime,
        EndTime:      query.EndTime,
    }
    return result, nil
}

func sumMetricValues(results []types.MetricDataResult) (float64, error) {
    var sum float64
    for _, result := range results {
        for _, val := range result.Values {
            sum += val
        }
    }
    return sum, nil
}