package metrics

import (
	"context"
	"fmt"
	"strconv"

	logsinsightsfetcher "github.com/dominikhei/aws-lambda-analyzer/sdk/internal/logsinsights"
	"github.com/dominikhei/aws-lambda-analyzer/sdk/internal/queries"
	sdktypes "github.com/dominikhei/aws-lambda-analyzer/sdk/types"
)

func GetColdStartRate(
	ctx context.Context,
	logsFetcher *logsinsightsfetcher.Fetcher,
	query sdktypes.FunctionQuery,
	period int32,
) (*sdktypes.ColdStartRateReturn, error) {
    

	results, err := logsFetcher.RunQuery(ctx, query, queries.LambdaColdStartRate)
    if err != nil {
        return nil, fmt.Errorf("run logs insights query: %w", err)
    }
    var coldStartRate float64
    if len(results) > 0 {
        val := results[0]["coldStartRate"]
        if val != "" {
            coldStartRate, err = strconv.ParseFloat(val, 32)
            if err != nil {
                return nil, fmt.Errorf("parse coldStartRate from logs: %w", err)
            }
        }
    }
	return &sdktypes.ColdStartRateReturn{
        ColdStartRate:  float32(coldStartRate),
        FunctionName: query.FunctionName,
        Qualifier:    query.Qualifier,
        StartTime:    query.StartTime,
        EndTime:      query.EndTime,
    }, nil
}