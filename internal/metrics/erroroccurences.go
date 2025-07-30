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

func GetErrorTypes(
	ctx context.Context,
	logsFetcher *logsinsightsfetcher.Fetcher,
	cwFetcher *cloudwatchfetcher.Fetcher,
	query sdktypes.FunctionQuery,
) (*sdktypes.ErrorTypesReturn, error) {

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
	queryString := fmt.Sprintf(queries.LambdaErrorTypesQueryWithVersion, escapedQualifier)
	results, err := logsFetcher.RunQuery(ctx, query, queryString)
	if err != nil {
		return nil, fmt.Errorf("run logs insights query: %w", err)
	}

	var stats []sdktypes.ErrorType
	for _, row := range results {
		category, ok := row["error_category"]
		if !ok || category == "" {
			category = "UnknownError"
		}
		countStr, ok := row["error_count"]
		if !ok {
			continue
		}
		count, err := strconv.ParseInt(countStr, 10, 64)
		if err != nil {
			fmt.Printf("warn: could not parse error_count %q: %v\n", countStr, err)
			continue
		}
		stats = append(stats, sdktypes.ErrorType{
			ErrorCategory: category,
			ErrorCount:    count,
		})
	}

	return &sdktypes.ErrorTypesReturn{
		Errors:       stats,
		FunctionName: query.FunctionName,
		Qualifier:    query.Qualifier,
		StartTime:    query.StartTime,
		EndTime:      query.EndTime,
	}, nil
}
