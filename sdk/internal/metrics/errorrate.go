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
    results, err := logsFetcher.RunQuery(ctx, query, queries.LambdaErrorRate)
    if err != nil {
        return nil, fmt.Errorf("run logs insights query: %w", err)
    }

    var total, errors int64
    for _, row := range results {
        if val, ok := row["totalInvocations"]; ok {
            total, _ = strconv.ParseInt(val, 10, 64)
        }
        if val, ok := row["errorInvocations"]; ok {
            errors, _ = strconv.ParseInt(val, 10, 64)
        }
    }

    if total == 0 {
        return &sdktypes.ErrorRateReturn{
            FunctionName: query.FunctionName,
            Qualifier:   query.Qualifier,
            StartTime:   query.StartTime,
            EndTime:     query.EndTime,
            ErrorRate:   0,
        }, nil
    }

    errorRate := float32(errors) / float32(total)

    return &sdktypes.ErrorRateReturn{
        FunctionName: query.FunctionName,
        Qualifier:   query.Qualifier,
        StartTime:   query.StartTime,
        EndTime:     query.EndTime,
        ErrorRate:   errorRate,
    }, nil
}
