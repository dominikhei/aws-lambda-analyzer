package metrics

import (
	"context"
	"fmt"

	sdkerrors "github.com/dominikhei/serverless-statistics/errors"
	cloudwatchfetcher "github.com/dominikhei/serverless-statistics/internal/cloudwatch"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

// GetErrorRate calculates the ratio of errors and total invoations
// of an AWS Lambda function over a specified time range and qualifier (version).
func GetErrorRate(
	ctx context.Context,
	cwFetcher *cloudwatchfetcher.Fetcher,
	query sdktypes.FunctionQuery,
) (*sdktypes.ErrorRateReturn, error) {
	invocationsResults, err := cwFetcher.FetchMetric(ctx, query, "Invocations", "Sum")
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

	errorsResults, err := cwFetcher.FetchMetric(ctx, query, "Errors", "Sum")
	if err != nil {
		return nil, fmt.Errorf("fetch errors metric: %w", err)
	}
	errorsSum, err := sumMetricValues(errorsResults)
	if err != nil {
		return nil, fmt.Errorf("parse errors metric data: %w", err)
	}

	errorRate := float64(errorsSum) / float64(invocationsSum)

	return &sdktypes.ErrorRateReturn{
		FunctionName: query.FunctionName,
		Qualifier:    query.Qualifier,
		StartTime:    query.StartTime,
		EndTime:      query.EndTime,
		ErrorRate:    errorRate,
	}, nil
}
