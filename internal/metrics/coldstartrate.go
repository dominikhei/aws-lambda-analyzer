package metrics

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	sdkerrors "github.com/dominikhei/serverless-statistics/errors"
	cloudwatchfetcher "github.com/dominikhei/serverless-statistics/internal/cloudwatch"
	logsinsightsfetcher "github.com/dominikhei/serverless-statistics/internal/logsinsights"
	"github.com/dominikhei/serverless-statistics/internal/queries"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

func GetColdStartRate(
	ctx context.Context,
	logsFetcher *logsinsightsfetcher.Fetcher,
	cwFetcher *cloudwatchfetcher.Fetcher,
	query sdktypes.FunctionQuery,
) (*sdktypes.ColdStartRateReturn, error) {

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

	escapedQualifier := strings.ReplaceAll(query.Qualifier, "$", "\\$")
	queryString := fmt.Sprintf(queries.LambdaColdStartRateWithVersion, escapedQualifier)
	results, err := logsFetcher.RunQuery(ctx, query, queryString)
	if err != nil {
		return nil, fmt.Errorf("run logs insights query: %w", err)
	}
	var coldStartRate float64
	totalStr := results[0]["totalInvocations"]
	coldStr := results[0]["coldStartLines"]

	if totalStr != "" && coldStr != "" {
		total, err1 := strconv.ParseFloat(totalStr, 32)
		cold, err2 := strconv.ParseFloat(coldStr, 32)
		if err1 != nil || err2 != nil || total == 0 {
			return nil, fmt.Errorf("invalid data from logs: total=%v, cold=%v", totalStr, coldStr)
		}
		coldStartRate = cold / total
	}
	return &sdktypes.ColdStartRateReturn{
		ColdStartRate: float32(coldStartRate),
		FunctionName:  query.FunctionName,
		Qualifier:     query.Qualifier,
		StartTime:     query.StartTime,
		EndTime:       query.EndTime,
	}, nil
}
