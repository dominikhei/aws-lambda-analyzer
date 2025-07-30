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

// This function does not use the duration metric from cloudwatch, as
// there is a risk of aggregating durations depending on the period.

// GetDurationStatistics calculates duration statistics of an AWS Lambda function
// over a specified time range and qualifier (version).
// It does not use the duration metric from cloudwatch, as
// there is a risk of aggregating durations depending on the period
func GetDurationStatistics(
	ctx context.Context,
	logsFetcher *logsinsightsfetcher.Fetcher,
	cwFetcher *cloudwatchfetcher.Fetcher,
	query sdktypes.FunctionQuery,
) (*sdktypes.DurationStatisticsReturn, error) {
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
	queryString := fmt.Sprintf(queries.LambdaDurationQueryWithVersion, escapedQualifier)
	results, err := logsFetcher.RunQuery(ctx, query, queryString)
	if err != nil {
		return nil, fmt.Errorf("run logs insights query: %w", err)
	}
	var durations []float64
	for _, row := range results {
		if valStr, ok := row["durationMs"]; ok {
			if val, err := strconv.ParseFloat(valStr, 64); err == nil {
				durations = append(durations, val)
			} else {
				fmt.Printf("warn: could not parse %q as float64: %v", valStr, err)
			}
		}
	}
	durationStats, err := utils.CalcSummaryStats(durations)
	if err != nil {
		return nil, fmt.Errorf("error calculating summary statistics: %w", err)
	}
	return &sdktypes.DurationStatisticsReturn{
		MinDuration:    durationStats.Min,
		MaxDuration:    durationStats.Max,
		MedianDuration: durationStats.Median,
		MeanDuration:   durationStats.Mean,
		P95Duration:    durationStats.P95,
		P99Duration:    durationStats.P99,
		Conf95Duration: durationStats.ConfInt95,
		FunctionName:   query.FunctionName,
		Qualifier:      query.Qualifier,
		StartTime:      query.StartTime,
		EndTime:        query.EndTime,
	}, nil
}
