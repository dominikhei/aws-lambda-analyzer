package metrics

import (
	"context"
	"fmt"
	cloudwatchfetcher "github.com/dominikhei/aws-lambda-analyzer/sdk/internal/cloudwatch"
	sdktypes "github.com/dominikhei/aws-lambda-analyzer/sdk/types"
	sdkerrors "github.com/dominikhei/aws-lambda-analyzer/sdk/errors"
)

func GetErrorRate(
	ctx context.Context,
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

	errorsResults, err := cwFetcher.FetchMetric(ctx, query, "Errors", "Sum", period)
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
		ErrorRate:    float32(errorRate),
	}, nil
}
