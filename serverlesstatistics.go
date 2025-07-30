package serverlessstatistics

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/dominikhei/serverless-statistics/internal/clientmanager"
	cloudwatchfetcher "github.com/dominikhei/serverless-statistics/internal/cloudwatch"
	logsinsightsfetcher "github.com/dominikhei/serverless-statistics/internal/logsinsights"
	"github.com/dominikhei/serverless-statistics/internal/metrics"
	"github.com/dominikhei/serverless-statistics/internal/utils"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

type ServerlessStats struct {
	cloudwatchFetcher *cloudwatchfetcher.Fetcher
	logsFetcher       *logsinsightsfetcher.Fetcher
	lambdaClient      *lambda.Client
}

func New(ctx context.Context, opts sdktypes.ConfigOptions) *ServerlessStats {
	clients, err := clientmanager.NewAWSClients(ctx, opts)
	if err != nil {
		log.Fatalf("failed to initialize AWS clients: %v", err)
	}

	return &ServerlessStats{
		cloudwatchFetcher: cloudwatchfetcher.New(clients),
		logsFetcher:       logsinsightsfetcher.New(clients),
		lambdaClient:      clients.LambdaClient,
	}
}

func (a *ServerlessStats) GetThrottleRate(
	ctx context.Context,
	functionName string,
	qualifier string,
	startTime, endTime time.Time,
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

	return metrics.GetThrottleRate(ctx, a.cloudwatchFetcher, query)
}

func (a *ServerlessStats) GetTimeoutRate(
	ctx context.Context,
	functionName string,
	qualifier string,
	startTime, endTime time.Time,
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

	return metrics.GetTimeoutRate(ctx, a.cloudwatchFetcher, a.logsFetcher, query)
}

func (a *ServerlessStats) GetColdStartRate(
	ctx context.Context,
	functionName string,
	qualifier string,
	startTime, endTime time.Time,
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

	return metrics.GetColdStartRate(ctx, a.logsFetcher, a.cloudwatchFetcher, query)
}

func (a *ServerlessStats) GetMaxMemoryUsageStatistics(
	ctx context.Context,
	functionName string,
	qualifier string,
	startTime, endTime time.Time,
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

	return metrics.GetMaxMemoryUsageStatistics(ctx, a.logsFetcher, a.cloudwatchFetcher, query)
}

func (a *ServerlessStats) GetErrorRate(
	ctx context.Context,
	functionName string,
	qualifier string,
	startTime, endTime time.Time,
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

	return metrics.GetErrorRate(ctx, a.cloudwatchFetcher, query)
}

func (a *ServerlessStats) GetErrorCategoryStatistics(
	ctx context.Context,
	functionName string,
	qualifier string,
	startTime, endTime time.Time,
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

	return metrics.GetErrorTypes(ctx, a.logsFetcher, a.cloudwatchFetcher, query)
}

func (a *ServerlessStats) GetDurationStatistics(
	ctx context.Context,
	functionName string,
	qualifier string,
	startTime, endTime time.Time,
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

	return metrics.GetDurationStatistics(ctx, a.logsFetcher, a.cloudwatchFetcher, query)
}
