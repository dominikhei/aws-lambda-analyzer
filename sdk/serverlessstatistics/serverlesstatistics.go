package serverlessstatistics

import (
	"context"
	"log"
	"time"

	"github.com/dominikhei/aws-lambda-analyzer/sdk/internal/clientmanager"
	"github.com/dominikhei/aws-lambda-analyzer/sdk/internal/cloudwatch"
	logsinsightsfetcher "github.com/dominikhei/aws-lambda-analyzer/sdk/internal/logsinsights"
	"github.com/dominikhei/aws-lambda-analyzer/sdk/internal/metrics"
	sdktypes "github.com/dominikhei/aws-lambda-analyzer/sdk/types"
)

type ServerlessStats struct {
    cloudwatchfetcher *cloudwatchfetcher.Fetcher
    logsFetcher *logsinsightsfetcher.Fetcher
}

func New(ctx context.Context, opts sdktypes.ConfigOptions) *ServerlessStats {
    clients, err := clientmanager.NewAWSClients(ctx, opts)
    if err != nil {
        log.Fatalf("failed to initialize AWS clients: %v", err)
    }

    return &ServerlessStats{
        cloudwatchfetcher: cloudwatchfetcher.New(clients),
        logsFetcher: logsinsightsfetcher.New(clients),
    }
}

func (a *ServerlessStats) GetThrottleRate(
    ctx context.Context,
    functionName string,
    qualifier string,
    startTime, endTime time.Time,
    period int32,
) (*sdktypes.ThrottleRateReturn, error) {
    query := sdktypes.FunctionQuery{
        FunctionName: functionName,
        Qualifier:    qualifier,
        StartTime:    startTime,
        EndTime:      endTime,
    }
    return metrics.GetThrottleRate(ctx, a.cloudwatchfetcher, query, period)
}

func (a *ServerlessStats) GetTimeoutRate(
    ctx context.Context,
    functionName string,
    qualifier string,
    startTime, endTime time.Time,
    period int32,
) (*sdktypes.TimeoutRateReturn, error) {
    query := sdktypes.FunctionQuery{
        FunctionName: functionName,
        Qualifier:    qualifier,
        StartTime:    startTime,
        EndTime:      endTime,
    }
    return metrics.GetTimeoutRate(ctx, a.cloudwatchfetcher, a.logsFetcher, query, period)
}

func (a *ServerlessStats) GetColdStartRate(
    ctx context.Context,
    functionARN string,
    qualifier string,
    startTime, endTime time.Time,
    period int32,    
) (*sdktypes.ColdStartRateReturn, error) {
    query := sdktypes.FunctionQuery{
        FunctionName: functionARN,
        Qualifier:    qualifier,
        StartTime:    startTime,
        EndTime:      endTime,
    }  
    return metrics.GetColdStartRate(ctx, a.logsFetcher, query, period)
}

func (a *ServerlessStats) GetMaxMemoryUsageStatistics(
    ctx context.Context,
    functionARN string,
    qualifier string,
    startTime, endTime time.Time,
    period int32,
) (*sdktypes.MemoryUsagePercentilesReturn, error) {
    query := sdktypes.FunctionQuery{
        FunctionName: functionARN,
        Qualifier:    qualifier,
        StartTime:    startTime,
        EndTime:      endTime,
    }
    return metrics.GetMaxMemoryUsageStatistics(ctx, a.logsFetcher, query, period)
}

func (a *ServerlessStats) GetErrorRate(
    ctx context.Context,
    functionARN string,
    qualifier string,
    startTime, endTime time.Time,
    period int32,
) (*sdktypes.ErrorRateReturn, error) {
    query := sdktypes.FunctionQuery{
        FunctionName: functionARN,
        Qualifier:    qualifier,
        StartTime:    startTime,
        EndTime:      endTime,
    }
    return metrics.GetErrorRate(ctx, a.logsFetcher, query, period)
}