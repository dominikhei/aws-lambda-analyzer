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
	"github.com/dominikhei/serverless-statistics/internal/utils"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

func GetMaxMemoryUsageStatistics(
	ctx context.Context,
	logsFetcher *logsinsightsfetcher.Fetcher,
	cwFetcher *cloudwatchfetcher.Fetcher,
	query sdktypes.FunctionQuery,
) (*sdktypes.MemoryUsagePercentilesReturn, error) {

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
	queryString := fmt.Sprintf(queries.LambdaMemoryUtilizationQueryWithVersion, escapedQualifier)
	results, err := logsFetcher.RunQuery(ctx, query, queryString)
	if err != nil {
		return nil, fmt.Errorf("run logs insights query: %w", err)
	}

	var ratios []float64
	for _, row := range results {
		if valStr, ok := row["memoryUtilizationRatio"]; ok {
			if val, err := strconv.ParseFloat(valStr, 64); err == nil {
				ratios = append(ratios, val)
			} else {
				fmt.Printf("warn: could not parse %q as float64: %v", valStr, err)
			}
		}
	}
	memoryStats, err := utils.CalcSummaryStats(ratios)
	if err != nil {
		return nil, fmt.Errorf("error calculating summary statistics: %w", err)
	}
	return &sdktypes.MemoryUsagePercentilesReturn{
		MinUsageRate:    memoryStats.Min,
		MaxUsageRate:    memoryStats.Max,
		MedianUsageRate: memoryStats.Median,
		MeanUsageRate:   memoryStats.Mean,
		P95UsageRate:    memoryStats.P95,
		P99UsageRate:    memoryStats.P99,
		Conf95UsageRate: memoryStats.ConfInt95,
		FunctionName:    query.FunctionName,
		Qualifier:       query.Qualifier,
		StartTime:       query.StartTime,
		EndTime:         query.EndTime,
	}, nil
}
