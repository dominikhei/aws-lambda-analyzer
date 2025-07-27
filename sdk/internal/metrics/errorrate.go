package metrics

import (
	"context"
	"fmt"
	"strconv"

	logsinsightsfetcher "github.com/dominikhei/aws-lambda-analyzer/sdk/internal/logsinsights"
	"github.com/dominikhei/aws-lambda-analyzer/sdk/internal/queries"
	sdktypes "github.com/dominikhei/aws-lambda-analyzer/sdk/types"
)

func GetErrorRate(
    ctx context.Context,
    logsFetcher *logsinsightsfetcher.Fetcher,
    query sdktypes.FunctionQuery,
	period int32,
) (*sdktypes.ErrorRateReturn, error) {
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
