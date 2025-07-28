package metrics

import (
	"context"
	"fmt"
	"strconv"

	cloudwatchfetcher "github.com/dominikhei/aws-lambda-analyzer/sdk/internal/cloudwatch"
	logsinsightsfetcher "github.com/dominikhei/aws-lambda-analyzer/sdk/internal/logsinsights"
	"github.com/dominikhei/aws-lambda-analyzer/sdk/internal/queries"
	sdktypes "github.com/dominikhei/aws-lambda-analyzer/sdk/types"
    sdkerrors "github.com/dominikhei/aws-lambda-analyzer/sdk/errors"
)

func GetErrorRate(
    ctx context.Context,
    logsFetcher *logsinsightsfetcher.Fetcher,
    cwFetcher *cloudwatchfetcher.Fetcher,
    query sdktypes.FunctionQuery,
	period int32,
) (*sdktypes.ErrorRateReturn, error) {

    invocationsResults, err := cwFetcher.FetchMetric(ctx, query, "Invocations", "Sum", period)
    if err != nil {
        return nil, fmt.Errorf("fetch invocations metric: %w", err)
    }
    invocationsSum, err := sumMetricValues(invocationsResults)
    if err != nil {
        return nil, fmt.Errorf("parse invocations metric data: %w", err)
    }
    if invocationsSum == 0 {
        return nil, sdkerrors.NewNoInvocationsError(query.FunctionName)
    }

    results, err := logsFetcher.RunQuery(ctx, query, queries.LambdaErrorCount)
    if err != nil {
        return nil, fmt.Errorf("fetch errors from logs insights: %w", err)
    }    
    var errorCount int
    if len(results) > 0 {
        val := results[0]["errorCount"]
        if val != "" {
            errorCount, err = strconv.Atoi(val)
            if err != nil {
                return nil, fmt.Errorf("parse errorCount from logs: %w", err)
            }
        }
    }

    results, err = logsFetcher.RunQuery(ctx, query, queries.LambdaUniqueRequests)
    if err != nil {
        return nil, fmt.Errorf("fetch errors from logs insights: %w", err)
    }    
    var invocationsCount int
    if len(results) > 0 {
        val := results[0]["invocationsCount"]
        if val != "" {
            invocationsCount, err = strconv.Atoi(val)
            if err != nil {
                return nil, fmt.Errorf("parse invocationsCount from logs: %w", err)
            }
        }
    }


    errorRate := float64(errorCount) / float64(invocationsCount)

    return &sdktypes.ErrorRateReturn{
        FunctionName: query.FunctionName,
        Qualifier:   query.Qualifier,
        StartTime:   query.StartTime,
        EndTime:     query.EndTime,
        ErrorRate:   float32(errorRate),
    }, nil
}
