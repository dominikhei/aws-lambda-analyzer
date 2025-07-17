package analyzer

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

type Analyzer struct {
    cloudwatchfetcher *cloudwatchfetcher.Fetcher
    logsFetcher *logsinsightsfetcher.Fetcher
}

func New(ctx context.Context, opts sdktypes.ConfigOptions) *Analyzer {
    clients, err := clientmanager.NewAWSClients(ctx, opts)
    if err != nil {
        log.Fatalf("failed to initialize AWS clients: %v", err)
    }

    return &Analyzer{
        cloudwatchfetcher: cloudwatchfetcher.New(clients),
        logsFetcher: logsinsightsfetcher.New(clients),
    }
}

func (a *Analyzer) GetThrottleRate(
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

func (a *Analyzer) GetTimeoutRate(
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
