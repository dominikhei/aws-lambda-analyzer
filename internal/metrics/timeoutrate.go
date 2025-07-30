package metrics

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	cloudwatchfetcher "github.com/dominikhei/serverless-statistics/internal/cloudwatch"
	logsinsightsfetcher "github.com/dominikhei/serverless-statistics/internal/logsinsights"
	"github.com/dominikhei/serverless-statistics/internal/queries"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

func GetTimeoutRate(
	ctx context.Context,
	cwFetcher *cloudwatchfetcher.Fetcher,
	logsFetcher *logsinsightsfetcher.Fetcher,
	query sdktypes.FunctionQuery,
) (*sdktypes.TimeoutRateReturn, error) {

	escapedQualifier := strings.ReplaceAll(query.Qualifier, "$", "\\$")
	queryString := fmt.Sprintf(queries.LambdaUniqueRequestsWithVersion, escapedQualifier)
	results, err := logsFetcher.RunQuery(ctx, query, queryString)
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
	escapedQualifier = strings.ReplaceAll(query.Qualifier, "$", "\\$")
	queryString = fmt.Sprintf(queries.LambdaTimeoutQueryWithVersion, escapedQualifier)
	results, err = logsFetcher.RunQuery(ctx, query, queryString)
	if err != nil {
		return nil, fmt.Errorf("run logs insights query: %w", err)
	}
	var timeoutCount int
	if len(results) > 0 {
		val := results[0]["timeoutCount"]
		if val != "" {
			timeoutCount, err = strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("parse timeoutCount from logs: %w", err)
			}
		}
	}
	timeoutRate := float64(timeoutCount) / float64(invocationsCount)

	return &sdktypes.TimeoutRateReturn{
		TimeoutRate:  timeoutRate,
		FunctionName: query.FunctionName,
		Qualifier:    query.Qualifier,
		StartTime:    query.StartTime,
		EndTime:      query.EndTime,
	}, nil
}
