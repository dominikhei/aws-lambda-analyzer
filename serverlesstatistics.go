// Copyright 2025 dominikhei
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// New initializes and returns a new instance of ServerlessStats.
// It sets up AWS service clients (CloudWatch, Lambda, and CloudWatch Logs Insights)
// using the provided configuration options. If client initialization fails,
// the function logs a fatal error and terminates the application.
//
// Example:
//
//	ctx := context.Background()
//	opts := types.ConfigOptions{
//		Region: "us-west-2",
//		// Credentials, profile, or other options...
//	}
//	stats := serverlessstatistics.New(ctx, opts)
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

// GetThrottleRate returns the rate of throttled invocations for a given
// Lambda function and version within the specified time range.
//
// Example:
//
//	throttleReturn, err := serverlessstatistics.GetThrottleRate("my-function", "v1", time.Now().Add(-1*time.Hour), time.Now())
//	if err != nil {
//		log.Fatalf("failed to get throttle rate: %v", err)
//	}
//	fmt.Printf("Throttle rate: %.2f%%\n", throttleReturn.ThrottleRate * 100)
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

// GetTimeoutRate returns the rate of timed-out invocations for a given
// Lambda function and version within the specified time range.
//
//	timeoutReturn, err := serverlessstatistics.GetTimeoutRate(ctx, "my-function", "v1", time.Now().Add(-1*time.Hour), time.Now())
//	if err != nil {
//		log.Fatalf("failed to get timeout rate: %v", err)
//	}
//	fmt.Printf("Timeout rate: %.2f%%\n", timeoutReturn.TimeoutRate * 100)
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

// GetColdStartRate returns the rate of cold start invocations for a given
// Lambda function and version within the specified time range.
//
// Example:
//
//	coldStartReturn, err := serverlessstatistics.GetColdStartRate(ctx, "my-function", "v1", time.Now().Add(-1*time.Hour), time.Now())
//	if err != nil {
//		log.Fatalf("failed to get cold start rate: %v", err)
//	}
//	fmt.Printf("Cold start rate: %.2f%%\n", coldStartReturn.ColdStartRate * 100)
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

// GetMaxMemoryUsageStatistics returns memory usage statistics for a given
// Lambda function and version within the specified time range.
// Example:
//
//	memoryReturn, err := serverlessstatistics.GetMaxMemoryUsageStatistics(ctx, "my-function", "v1", time.Now().Add(-1*time.Hour), time.Now())
//	if err != nil {
//		log.Fatalf("failed to get memory stats: %v", err)
//	}
//	if memoryReturn.P99UsageRate != nil {
//		fmt.Printf("P99 memory usage: %.2f MB\n", memoryReturn.P99UsageRate)
//	}
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

// GetErrorRate returns the rate of error invocations for a given
// Lambda function and version within the specified time range.
//
// Example:
//
//	errorReturn, err := serverlessstatistics.GetErrorRate(ctx, "my-function", "v1", time.Now().Add(-1*time.Hour), time.Now())
//	if err != nil {
//		log.Fatalf("failed to get error rate: %v", err)
//	}
//	fmt.Printf("Error rate: %.2f%%\n", errorReturn.ErrorRate * 100)
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

// GetErrorCategoryStatistics returns the distribution of error categories for a given
// Lambda function and version within the specified time range.
//
// Example:
//
//	errorCategoryReturn, err := serverlessstatistics.GetErrorCategoryStatistics(ctx, "my-function", "v1", time.Now().Add(-1*time.Hour), time.Now())
//	if err != nil {
//		log.Fatalf("failed to get error categories: %v", err)
//	}
//	for _, errType := range errorCategoryReturn.Errors {
//		fmt.Printf("Category: %s, Count: %d\n", errType.ErrorCategory, errType.ErrorCount)
//	}
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

// GetDurationStatistics returns duration statistics for a given
// Lambda function and version within the specified time range.
//
// Example:
//
//	durationReturn, err := serverlessstatistics.GetDurationStatistics(ctx, "my-function", "v1", time.Now().Add(-1*time.Hour), time.Now())
//	if err != nil {
//		log.Fatalf("failed to get duration statistics: %v", err)
//	}
//	if durationReturn.P99Duration != nil {
//		fmt.Printf("P99 duration: %.2f MB\n", durationReturn.P99Duration)
//	}
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
