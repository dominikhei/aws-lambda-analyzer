package serverlessstatistics

import (
	"context"
	"log"
	"time"
    "fmt"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/dominikhei/aws-lambda-analyzer/sdk/internal/clientmanager"
	"github.com/dominikhei/aws-lambda-analyzer/sdk/internal/cloudwatch"
	logsinsightsfetcher "github.com/dominikhei/aws-lambda-analyzer/sdk/internal/logsinsights"
	"github.com/dominikhei/aws-lambda-analyzer/sdk/internal/metrics"
	sdktypes "github.com/dominikhei/aws-lambda-analyzer/sdk/types"
    "github.com/dominikhei/aws-lambda-analyzer/sdk/internal/utils"
)

type ServerlessStats struct {
    cloudwatchFetcher *cloudwatchfetcher.Fetcher
    logsFetcher *logsinsightsfetcher.Fetcher
    lambdaClient      *lambda.Client
}

func New(ctx context.Context, opts sdktypes.ConfigOptions) *ServerlessStats {
    clients, err := clientmanager.NewAWSClients(ctx, opts)
    if err != nil {
        log.Fatalf("failed to initialize AWS clients: %v", err)
    }

    return &ServerlessStats{
        cloudwatchFetcher: cloudwatchfetcher.New(clients),
        logsFetcher: logsinsightsfetcher.New(clients),
        lambdaClient: clients.LambdaClient,
    }
}

func (a *ServerlessStats) GetThrottleRate(
    ctx context.Context,
    functionName string,
    qualifier string,
    startTime, endTime time.Time,
    period int32,
) (*sdktypes.ThrottleRateReturn, error) {
    if qualifier == "" {
        qualifier = "$LATEST"
    }
    query := sdktypes.FunctionQuery{
        FunctionName: functionName,
        Qualifier:    qualifier,
        StartTime:    startTime,
        EndTime:      endTime,
    }

    exists, err := utils.FunctionExists(ctx, a.lambdaClient, functionName)
    if err != nil {
        return nil, fmt.Errorf("checking if function exists: %w", err)
    }
    if !exists {
        return nil, fmt.Errorf("lambda function %q does not exist", functionName)
    }

    exists, err = utils.QualifierExists(ctx, a.lambdaClient, functionName, qualifier)
    if err != nil {
        return nil, fmt.Errorf("checking if qualifier exists: %w", err)
    }
    if !exists {
        return nil, fmt.Errorf("qualifier %q does not exist", qualifier)
    }

    return metrics.GetThrottleRate(ctx, a.cloudwatchFetcher, query, period)
}

func (a *ServerlessStats) GetTimeoutRate(
    ctx context.Context,
    functionName string,
    qualifier string,
    startTime, endTime time.Time,
    period int32,
) (*sdktypes.TimeoutRateReturn, error) {
    if qualifier == "" {
        qualifier = "$LATEST"
    }
    query := sdktypes.FunctionQuery{
        FunctionName: functionName,
        Qualifier:    qualifier,
        StartTime:    startTime,
        EndTime:      endTime,
    }

    exists, err := utils.FunctionExists(ctx, a.lambdaClient, functionName)
    if err != nil {
        return nil, fmt.Errorf("checking if function exists: %w", err)
    }
    if !exists {
        return nil, fmt.Errorf("lambda function %q does not exist", functionName)
    }

    exists, err = utils.QualifierExists(ctx, a.lambdaClient, functionName, qualifier)
    if err != nil {
        return nil, fmt.Errorf("checking if qualifier exists: %w", err)
    }
    if !exists {
        return nil, fmt.Errorf("qualifier %q does not exist", qualifier)
    }

    return metrics.GetTimeoutRate(ctx, a.cloudwatchFetcher, a.logsFetcher, query, period)
}

func (a *ServerlessStats) GetColdStartRate(
    ctx context.Context,
    functionName string,
    qualifier string,
    startTime, endTime time.Time,
    period int32,    
) (*sdktypes.ColdStartRateReturn, error) {
    if qualifier == "" {
        qualifier = "$LATEST"
    }
    query := sdktypes.FunctionQuery{
        FunctionName: functionName,
        Qualifier:    qualifier,
        StartTime:    startTime,
        EndTime:      endTime,
    }  

    exists, err := utils.FunctionExists(ctx, a.lambdaClient, functionName)
    if err != nil {
        return nil, fmt.Errorf("checking if function exists: %w", err)
    }
    if !exists {
        return nil, fmt.Errorf("lambda function %q does not exist", functionName)
    }

    exists, err = utils.QualifierExists(ctx, a.lambdaClient, functionName, qualifier)
    if err != nil {
        return nil, fmt.Errorf("checking if qualifier exists: %w", err)
    }
    if !exists {
        return nil, fmt.Errorf("qualifier %q does not exist", qualifier)
    }

    return metrics.GetColdStartRate(ctx, a.logsFetcher, a.cloudwatchFetcher, query, period)
}

func (a *ServerlessStats) GetMaxMemoryUsageStatistics(
    ctx context.Context,
    functionName string,
    qualifier string,
    startTime, endTime time.Time,
    period int32,
) (*sdktypes.MemoryUsagePercentilesReturn, error) {
    if qualifier == "" {
        qualifier = "$LATEST"
    }
    query := sdktypes.FunctionQuery{
        FunctionName: functionName,
        Qualifier:    qualifier,
        StartTime:    startTime,
        EndTime:      endTime,
    }

    exists, err := utils.FunctionExists(ctx, a.lambdaClient, functionName)
    if err != nil {
        return nil, fmt.Errorf("checking if function exists: %w", err)
    }
    if !exists {
        return nil, fmt.Errorf("lambda function %q does not exist", functionName)
    }

    exists, err = utils.QualifierExists(ctx, a.lambdaClient, functionName, qualifier)
    if err != nil {
        return nil, fmt.Errorf("checking if qualifier exists: %w", err)
    }
    if !exists {
        return nil, fmt.Errorf("qualifier %q does not exist", qualifier)
    }

    return metrics.GetMaxMemoryUsageStatistics(ctx, a.logsFetcher, a.cloudwatchFetcher, query, period)
}

func (a *ServerlessStats) GetErrorRate(
    ctx context.Context,
    functionName string,
    qualifier string,
    startTime, endTime time.Time,
    period int32,
) (*sdktypes.ErrorRateReturn, error) {
    if qualifier == "" {
        qualifier = "$LATEST"
    }
    query := sdktypes.FunctionQuery{
        FunctionName: functionName,
        Qualifier:    qualifier,
        StartTime:    startTime,
        EndTime:      endTime,
    }

    exists, err := utils.FunctionExists(ctx, a.lambdaClient, functionName)
    if err != nil {
        return nil, fmt.Errorf("checking if function exists: %w", err)
    }
    if !exists {
        return nil, fmt.Errorf("lambda function %q does not exist", functionName)
    }

    exists, err = utils.QualifierExists(ctx, a.lambdaClient, functionName, qualifier)
    if err != nil {
        return nil, fmt.Errorf("checking if qualifier exists: %w", err)
    }
    if !exists {
        return nil, fmt.Errorf("qualifier %q does not exist", qualifier)
    }

    return metrics.GetErrorRate(ctx, a.cloudwatchFetcher, query, period)
}

func (a *ServerlessStats) GetErrorCategoryStatistics(
    ctx context.Context,
    functionName string,
    qualifier string,
    startTime, endTime time.Time,
    period int32,
) (*sdktypes.ErrorTypesReturn, error) {
    if qualifier == "" {
        qualifier = "$LATEST"
    }
    query := sdktypes.FunctionQuery{
        FunctionName: functionName,
        Qualifier:    qualifier,
        StartTime:    startTime,
        EndTime:      endTime,
    }

    exists, err := utils.FunctionExists(ctx, a.lambdaClient, functionName)
    if err != nil {
        return nil, fmt.Errorf("checking if function exists: %w", err)
    }
    if !exists {
        return nil, fmt.Errorf("lambda function %q does not exist", functionName)
    }

    exists, err = utils.QualifierExists(ctx, a.lambdaClient, functionName, qualifier)
    if err != nil {
        return nil, fmt.Errorf("checking if qualifier exists: %w", err)
    }
    if !exists {
        return nil, fmt.Errorf("qualifier %q does not exist", qualifier)
    }

    return metrics.GetErrorTypes(ctx, a.logsFetcher, a.cloudwatchFetcher, query, period)
}

func (a *ServerlessStats) GetDurationStatistics(
    ctx context.Context,
    functionName string,
    qualifier string,
    startTime, endTime time.Time,
    period int32,
) (*sdktypes.DurationStatisticsReturn, error) {
    if qualifier == "" {
        qualifier = "$LATEST"
    }
    query := sdktypes.FunctionQuery{
        FunctionName: functionName,
        Qualifier:    qualifier,
        StartTime:    startTime,
        EndTime:      endTime,
    }

    exists, err := utils.FunctionExists(ctx, a.lambdaClient, functionName)
    if err != nil {
        return nil, fmt.Errorf("checking if function exists: %w", err)
    }
    if !exists {
        return nil, fmt.Errorf("lambda function %q does not exist", functionName)
    }

    exists, err = utils.QualifierExists(ctx, a.lambdaClient, functionName, qualifier)
    if err != nil {
        return nil, fmt.Errorf("checking if qualifier exists: %w", err)
    }
    if !exists {
        return nil, fmt.Errorf("qualifier %q does not exist", qualifier)
    }

    return metrics.GetDurationStatistics(ctx, a.logsFetcher, a.cloudwatchFetcher, query, period)
}