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

func GetColdStartRate(
	ctx context.Context,
	logsFetcher *logsinsightsfetcher.Fetcher,
    cwFetcher *cloudwatchfetcher.Fetcher,
	query sdktypes.FunctionQuery,
	period int32,
) (*sdktypes.ColdStartRateReturn, error) {

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

	results, err := logsFetcher.RunQuery(ctx, query, queries.LambdaColdStartRate)
    if err != nil {
        return nil, fmt.Errorf("run logs insights query: %w", err)
    }
    var coldStartRate float64
    if len(results) > 0 {
        val := results[0]["coldStartRate"]
        if val != "" {
            coldStartRate, err = strconv.ParseFloat(val, 32)
            if err != nil {
                return nil, fmt.Errorf("parse coldStartRate from logs: %w", err)
            }
        }
    }
	return &sdktypes.ColdStartRateReturn{
        ColdStartRate:  float32(coldStartRate),
        FunctionName: query.FunctionName,
        Qualifier:    query.Qualifier,
        StartTime:    query.StartTime,
        EndTime:      query.EndTime,
    }, nil
}