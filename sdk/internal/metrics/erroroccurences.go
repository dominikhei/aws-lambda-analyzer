package metrics

import (
	"context"
	"fmt"
	"strconv"

	logsinsightsfetcher "github.com/dominikhei/aws-lambda-analyzer/sdk/internal/logsinsights"
	"github.com/dominikhei/aws-lambda-analyzer/sdk/internal/queries"
	sdktypes "github.com/dominikhei/aws-lambda-analyzer/sdk/types"
)

func GetErrorTypes(
	ctx context.Context,
	logsFetcher *logsinsightsfetcher.Fetcher,
	query sdktypes.FunctionQuery,
	period int32,
) (*sdktypes.ErrorTypesReturn, error) {

	results, err := logsFetcher.RunQuery(ctx, query, queries.LambdaErrorTypesQuery)
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