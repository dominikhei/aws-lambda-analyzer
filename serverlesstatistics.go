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

// Package serverlessstatistics provides a comprehensive client for fetching
// detailed AWS Lambda performance and reliability metrics.
//
// It simplifies access to key Lambda statistics such as cold start rates,
// invocation durations, error rates and categories, memory usage percentiles,
// throttle and timeout rates, and waste ratios.
//
// This package wraps AWS CloudWatch, Lambda, and CloudWatch Logs Insights
// clients to gather and aggregate metrics over configurable time intervals.
//
// Typical usage involves initializing a ServerlessStats instance with your
// AWS configuration options and then querying for metrics for a specific
// Lambda function and version. Aliases are not supported.
//
// Example:
//
//	ctx := context.Background()
//	opts := types.ConfigOptions{
//		Region: "us-west-2",
//		// Additional AWS config options like credentials or profile can be set here
//	}
//	stats := serverlessstatistics.New(ctx, opts)
//
//	// Fetch cold start rate over the last 24 hours
//	startTime := time.Now().Add(-24 * time.Hour)
//	endTime := time.Now()
//	return, err := stats.GetColdStartRate(ctx, "myLambdaFunction", "$LATEST", startTime, endTime)
//	if err != nil {
//	    log.Fatalf("Failed to get cold start rate: %v", err)
//	}
//	fmt.Printf("Cold start rate: %f\n", return.ColdStartRate)
//
// The package is designed to aid developers and SRE teams in monitoring
// and optimizing the performance of their serverless applications by
// providing actionable Lambda function insights programmatically.
package serverlessstatistics

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/dominikhei/serverless-statistics/internal/cache"
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
	invocationsCache  *cache.Cache
}

// ServerlessStats holds clients and caches to fetch AWS Lambda statistics.
//
// New initializes and returns a new instance of ServerlessStats.
// It sets up AWS service clients (CloudWatch, Lambda, and CloudWatch Logs Insights)
// using the provided configuration options.
//
// If client initialization fails, it logs a fatal error and terminates the application.
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
		log.Fatalf("failed to initialize clients: %v", err)
	}

	return &ServerlessStats{
		cloudwatchFetcher: cloudwatchfetcher.New(clients),
		logsFetcher:       logsinsightsfetcher.New(clients),
		lambdaClient:      clients.LambdaClient,
		invocationsCache:  cache.NewCache(),
	}
}

// GetThrottleRate returns the throttle rate (i.e., the proportion of throttled invocations)
// for a given AWS Lambda function and version within the specified time range.
//
// Input Parameters:
//   - ctx: Context used for cancellation and timeouts.
//   - functionName: The name of the Lambda function to analyze.
//   - version: (Optional) Version of the Lambda function. If empty, defaults to "$LATEST".
//   - startTime: The beginning of the time window to analyze.
//   - endTime: The end of the time window to analyze.
//
// Returns:
//   - *sdktypes.ThrottleRateReturn: Struct containing the calculated throttle rate (as a float64 between 0 and 1).
//   - error: Non-nil if the function or version doesn't exist, or if the underlying CloudWatch fetch fails.
//
// Example:
//
//	throttleReturn, err := serverlessstatistics.GetThrottleRate(ctx, "my-function", "v1", time.Now().Add(-1*time.Hour), time.Now())
//	if err != nil {
//		log.Fatalf("failed to get throttle rate: %v", err)
//	}
//	fmt.Printf("Throttle rate: %.2f%%\n", throttleReturn.ThrottleRate * 100)
func (a *ServerlessStats) GetThrottleRate(
	ctx context.Context,
	functionName string,
	version string,
	startTime, endTime time.Time,
) (*sdktypes.ThrottleRateReturn, error) {
	if version == "" {
		version = "$LATEST"
	}
	query := sdktypes.FunctionQuery{
		FunctionName: functionName,
		Qualifier:    version,
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

	exists, err = utils.QualifierExists(ctx, a.lambdaClient, functionName, version)
	if err != nil {
		return nil, fmt.Errorf("checking if version exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("version %q does not exist", version)
	}

	return metrics.GetThrottleRate(ctx, a.cloudwatchFetcher, a.invocationsCache, query)
}

// GetTimeoutRate returns the timeout rate (i.e., the proportion of Lambda function
// invocations that timed out) for a given AWS Lambda function and version
// within the specified time range.
//
// Input Parameters:
//   - ctx: Context for cancellation, deadlines, and timeouts.
//   - functionName: The name of the AWS Lambda function to analyze.
//   - version: (Optional) Lambda version. If empty, defaults to "$LATEST".
//   - startTime: Start of the analysis time window (must be before endTime).
//   - endTime: End of the analysis time window (typically time.Now()).
//
// Returns:
//   - *sdktypes.TimeoutRateReturn: Struct containing the timeout rate as a float64 (e.g., 0.12 = 12%).
//   - error: Non-nil if the function or version doesn't exist, or if CloudWatch/Logs queries fail.
//
// Example:
//
//	timeoutReturn, err := serverlessstatistics.GetTimeoutRate(ctx, "my-function", "v1", time.Now().Add(-1*time.Hour), time.Now())
//	if err != nil {
//		log.Fatalf("failed to get timeout rate: %v", err)
//	}
//	fmt.Printf("Timeout rate: %.2f%%\n", timeoutReturn.TimeoutRate * 100)
func (a *ServerlessStats) GetTimeoutRate(
	ctx context.Context,
	functionName string,
	version string,
	startTime, endTime time.Time,
) (*sdktypes.TimeoutRateReturn, error) {
	if version == "" {
		version = "$LATEST"
	}
	query := sdktypes.FunctionQuery{
		FunctionName: functionName,
		Qualifier:    version,
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

	exists, err = utils.QualifierExists(ctx, a.lambdaClient, functionName, version)
	if err != nil {
		return nil, fmt.Errorf("checking if version exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("version %q does not exist", version)
	}

	return metrics.GetTimeoutRate(ctx, a.cloudwatchFetcher, a.logsFetcher, a.invocationsCache, query)
}

// GetColdStartRate returns the cold start rate for a given AWS Lambda function and version
// within the specified time range. A cold start is identified by the presence of an `initDuration` field
// in the invocation logs.
//
// Input Parameters:
//   - ctx: Context for timeout and cancellation control.
//   - functionName: The name of the AWS Lambda function to analyze.
//   - version: (Optional) Lambda version. If empty, defaults to "$LATEST".
//   - startTime: Start of the time window to analyze (should be within log retention).
//   - endTime: End of the time window to analyze (usually time.Now()).
//
// Returns:
//   - *sdktypes.ColdStartRateReturn: Struct containing the cold start rate as a float64 (e.g., 0.18 = 18%).
//   - error: Returned if the function or version does not exist or if log/metric queries fail.
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
	version string,
	startTime, endTime time.Time,
) (*sdktypes.ColdStartRateReturn, error) {
	if version == "" {
		version = "$LATEST"
	}
	query := sdktypes.FunctionQuery{
		FunctionName: functionName,
		Qualifier:    version,
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

	exists, err = utils.QualifierExists(ctx, a.lambdaClient, functionName, version)
	if err != nil {
		return nil, fmt.Errorf("checking if version exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("version %q does not exist", version)
	}

	return metrics.GetColdStartRate(ctx, a.logsFetcher, a.cloudwatchFetcher, a.invocationsCache, query)
}

// GetMaxMemoryUsageStatistics returns memory usage percentiles for a given AWS Lambda function
// and version within the specified time range. Memory usage is calculated as the
// ratio of maximum used memory to allocated memory for each invocation, and then aggregated
// into various percentiles.
//
// Input Parameters:
//   - ctx: Context for timeout and cancellation handling.
//   - functionName: The name of the AWS Lambda function to analyze.
//   - version: (Optional) Lambda version. If empty, defaults to "$LATEST".
//   - startTime: Start of the time window for analysis (should be within log retention).
//   - endTime: End of the time window for analysis (typically time.Now()).
//
// Returns:
//   - *sdktypes.MemoryUsagePercentilesReturn: Struct containing memory usage percentiles,
//     such as Min, Max, Median, Mean, P95, P99, and 95% Confidence Interval. Percentages are expressed as float64 (e.g., 0.73 = 73%).
//   - error: Returned if the function or version does not exist, or if log queries fail.
//
// Example:
//
//	memoryReturn, err := serverlessstatistics.GetMaxMemoryUsageStatistics(ctx, "my-function", "v1", time.Now().Add(-1*time.Hour), time.Now())
//	if err != nil {
//		log.Fatalf("failed to get memory stats: %v", err)
//	}
//	if memoryReturn.P99UsageRate != nil {
//		fmt.Printf("P99 memory usage: %.2f%%\n", *memoryReturn.P99UsageRate * 100)
//	}
func (a *ServerlessStats) GetMaxMemoryUsageStatistics(
	ctx context.Context,
	functionName string,
	version string,
	startTime, endTime time.Time,
) (*sdktypes.MemoryUsagePercentilesReturn, error) {
	if version == "" {
		version = "$LATEST"
	}
	query := sdktypes.FunctionQuery{
		FunctionName: functionName,
		Qualifier:    version,
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

	exists, err = utils.QualifierExists(ctx, a.lambdaClient, functionName, version)
	if err != nil {
		return nil, fmt.Errorf("checking if version exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("version %q does not exist", version)
	}

	return metrics.GetMaxMemoryUsageStatistics(ctx, a.logsFetcher, a.cloudwatchFetcher, a.invocationsCache, query)
}

// GetErrorRate returns the error rate for a given AWS Lambda function and version
// within the specified time range. An error is identified by log entries containing the "[ERROR]" tag.
//
// Input Parameters:
//   - ctx: Context for timeout and cancellation control.
//   - functionName: The name of the AWS Lambda function to analyze.
//   - version: (Optional) Lambda version. If empty, defaults to "$LATEST".
//   - startTime: Start of the time window to analyze (should be within log retention).
//   - endTime: End of the time window to analyze (typically time.Now()).
//
// Returns:
//   - *sdktypes.ErrorRateReturn: Struct containing the error rate as a float64 (e.g., 0.07 = 7%).
//   - error: Returned if the function or version does not exist, or if metric/log queries fail.
//
// Notes:
//   - This analysis is based on CloudWatch Logs Insights. It detects errors based on
//     presence of “[ERROR]” in logs tied to individual `requestId`s.
//   - Function timeouts and throttles are not included as “errors”.
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
	version string,
	startTime, endTime time.Time,
) (*sdktypes.ErrorRateReturn, error) {
	if version == "" {
		version = "$LATEST"
	}
	query := sdktypes.FunctionQuery{
		FunctionName: functionName,
		Qualifier:    version,
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

	exists, err = utils.QualifierExists(ctx, a.lambdaClient, functionName, version)
	if err != nil {
		return nil, fmt.Errorf("checking if version exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("version %q does not exist", version)
	}

	return metrics.GetErrorRate(ctx, a.cloudwatchFetcher, a.invocationsCache, query)
}

// GetErrorCategoryStatistics returns a categorized breakdown of errors for a given
// AWS Lambda function and version within the specified time range.
// Each error is grouped by its semantic type, extracted from log messages containing "[ERROR]".
//
// Input Parameters:
//   - ctx: Context for timeout and cancellation handling.
//   - functionName: The name of the AWS Lambda function to analyze.
//   - version: (Optional) Lambda version. If empty, defaults to "$LATEST".
//   - startTime: Start of the time window to analyze (must precede endTime and be within log retention).
//   - endTime: End of the time window to analyze (usually time.Now()).
//
// Returns:
//   - *sdktypes.ErrorTypesReturn: Struct containing a slice of error categories and their occurrence counts.
//   - error: Returned if the function or version does not exist, or if log queries fail.
//
// Notes:
//   - The grouping is based on log lines containing "[ERROR]", and the error type is extracted
//     semantically (e.g., "[ERROR] ImportError: ..." → `ImportError`).
//   - Timeouts are **not** classified as errors.
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
	version string,
	startTime, endTime time.Time,
) (*sdktypes.ErrorTypesReturn, error) {
	if version == "" {
		version = "$LATEST"
	}
	query := sdktypes.FunctionQuery{
		FunctionName: functionName,
		Qualifier:    version,
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

	exists, err = utils.QualifierExists(ctx, a.lambdaClient, functionName, version)
	if err != nil {
		return nil, fmt.Errorf("checking if version exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("version %q does not exist", version)
	}

	return metrics.GetErrorTypes(ctx, a.logsFetcher, a.cloudwatchFetcher, a.invocationsCache, query)
}

// GetDurationStatistics returns execution duration percentiles for a given AWS Lambda function
// and version within the specified time range. The duration refers to the time spent
// running the handler code (excluding init and billing overhead).
//
// Input Parameters:
//   - ctx: Context for timeout and cancellation control.
//   - functionName: The name of the AWS Lambda function to analyze.
//   - version: (Optional) Lambda version. If empty, defaults to "$LATEST".
//   - startTime: Start of the time window for analysis (must precede endTime).
//   - endTime: End of the time window for analysis (typically time.Now()).
//
// Returns:
//   - *sdktypes.DurationStatisticsReturn: Struct containing execution duration statistics such as
//     Min, Max, Median, Mean, P95, P99, and 95% Confidence Interval. Units are in milliseconds.
//   - error: Returned if the function or version does not exist, or if the underlying log/metric query fails.
//
// Notes:
//   - Durations are extracted from CloudWatch Logs (`REPORT` lines).
//   - Billing duration and cold start time are excluded; only handler execution is analyzed.
//   - Percentiles requiring a minimum number of invocations (e.g., P95, P99, CI) may be `nil`.
//
// Example:
//
//	durationReturn, err := serverlessstatistics.GetDurationStatistics(ctx, "my-function", "v1", time.Now().Add(-1*time.Hour), time.Now())
//	if err != nil {
//		log.Fatalf("failed to get duration statistics: %v", err)
//	}
//	if durationReturn.P99Duration != nil {
//		fmt.Printf("P99 duration: %.2f ms\n", *durationReturn.P99Duration)
//	}
func (a *ServerlessStats) GetDurationStatistics(
	ctx context.Context,
	functionName string,
	version string,
	startTime, endTime time.Time,
) (*sdktypes.DurationStatisticsReturn, error) {
	if version == "" {
		version = "$LATEST"
	}
	query := sdktypes.FunctionQuery{
		FunctionName: functionName,
		Qualifier:    version,
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

	exists, err = utils.QualifierExists(ctx, a.lambdaClient, functionName, version)
	if err != nil {
		return nil, fmt.Errorf("checking if version exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("version %q does not exist", version)
	}

	return metrics.GetDurationStatistics(ctx, a.logsFetcher, a.cloudwatchFetcher, a.invocationsCache, query)
}

// GetWasteRatio returns the ratio of billed duration that was not used by the handler execution
// for a given AWS Lambda function and version within the specified time range.
//
// The waste ratio quantifies the inefficiency of function executions in terms of over-allocated
// billing time (e.g., rounding up to the nearest 1 ms or 100 ms) compared to actual handler duration.
//
// Input Parameters:
//   - ctx: Context for timeout and cancellation control.
//   - functionName: The name of the AWS Lambda function to analyze.
//   - version: (Optional) Lambda version. If empty, defaults to "$LATEST".
//   - startTime: Start of the time window for analysis.
//   - endTime: End of the time window for analysis.
//
// Returns:
//   - *sdktypes.WasteRatioReturn: Struct containing the average waste ratio (0.0–1.0),
//     as well as optional breakdowns or supporting statistics.
//   - error: Returned if the function or version does not exist, or if metric/log retrieval fails.
//
// Notes:
//   - Waste ratio = (billed duration − actual duration) / billed duration.
//   - A waste ratio of 0.00 means no overhead; 0.25 means 25% of billed time was unused.
//
// Example:
//
//	wasteReturn, err := serverlessstatistics.GetWasteRatio(ctx, "my-function", "v1", time.Now().Add(-1*time.Hour), time.Now())
//	if err != nil {
//		log.Fatalf("failed to get waste ratio: %v", err)
//	}
//	if wasteReturn.WasteRatio != nil {
//		fmt.Printf("Waste ratio: %.2f%%\n", *wasteReturn.WasteRatio * 100)
//	}
func (a *ServerlessStats) GetWasteRatio(
	ctx context.Context,
	functionName string,
	version string,
	startTime, endTime time.Time,
) (*sdktypes.WasteRatioReturn, error) {
	if version == "" {
		version = "$LATEST"
	}
	query := sdktypes.FunctionQuery{
		FunctionName: functionName,
		Qualifier:    version,
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

	exists, err = utils.QualifierExists(ctx, a.lambdaClient, functionName, version)
	if err != nil {
		return nil, fmt.Errorf("checking if version exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("version %q does not exist", version)
	}

	return metrics.GetWasteRatio(ctx, a.cloudwatchFetcher, a.logsFetcher, a.invocationsCache, query)
}

// GetColdStartDurationStatistics returns statistics on cold start durations for a given
// AWS Lambda function and version within the specified time range.
//
// Cold start duration measures the additional latency incurred when Lambda initializes
// a new execution environment before invoking the function handler.
//
// Input Parameters:
//   - ctx: Context for cancellation and timeout.
//   - functionName: Name of the Lambda function to analyze.
//   - version: (Optional) Lambda version. Defaults to "$LATEST" if empty.
//   - startTime: Start timestamp for the analysis window.
//   - endTime: End timestamp for the analysis window.
//
// Returns:
//   - *sdktypes.ColdStartDurationStatisticsReturn: Struct containing percentiles (e.g., P99) of cold start durations.
//   - error: If the function or version does not exist, or if metrics/log retrieval fails.
//
// Example:
//
//	durationReturn, err := serverlessstatistics.GetColdStartDurationStatistics(ctx, "my-function", "v1", time.Now().Add(-1*time.Hour), time.Now())
//	if err != nil {
//		log.Fatalf("failed to get cold start duration statistics: %v", err)
//	}
//	if durationReturn.P99ColdStartDuration != nil {
//		fmt.Printf("P99 cold start duration: %.2f ms\n", *durationReturn.P99ColdStartDuration)
//	}
func (a *ServerlessStats) GetColdStartDurationStatistics(
	ctx context.Context,
	functionName string,
	version string,
	startTime, endTime time.Time,
) (*sdktypes.ColdStartDurationStatisticsReturn, error) {
	if version == "" {
		version = "$LATEST"
	}
	query := sdktypes.FunctionQuery{
		FunctionName: functionName,
		Qualifier:    version,
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

	exists, err = utils.QualifierExists(ctx, a.lambdaClient, functionName, version)
	if err != nil {
		return nil, fmt.Errorf("checking if version exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("version %q does not exist", version)
	}

	return metrics.GetColdStartDurationStatistics(ctx, a.logsFetcher, a.cloudwatchFetcher, a.invocationsCache, query)
}

// GetFunctionConfiguration returns the configuration details for a given
// AWS Lambda function and version.
//
// This includes metadata such as memory size, timeout, runtime, environment variables,
// and other configuration parameters.
//
// Input Parameters:
//   - ctx: Context for cancellation and timeout.
//   - functionName: Name of the Lambda function to retrieve configuration for.
//   - version: (Optional) Lambda version. Defaults to "$LATEST" if empty.
//
// Returns:
//   - *sdktypes.BaseStatisticsReturn: Struct containing the function's configuration details.
//   - error: If the function or version does not exist or retrieval fails.
//
// Example:
//
//	configs, err := serverlessstatistics.GetFunctionConfiguration(ctx, "my-function", "v1")
//	if err != nil {
//		log.Fatalf("failed to get function configuration: %v", err)
//	}
//	fmt.Printf("Memory size: %d MB\n", configs.MemorySize)
func (a *ServerlessStats) GetFunctionConfiguration(
	ctx context.Context,
	functionName string,
	version string,
) (*sdktypes.BaseStatisticsReturn, error) {
	if version == "" {
		version = "$LATEST"
	}
	query := sdktypes.FunctionQuery{
		FunctionName: functionName,
		Qualifier:    version,
	}

	exists, err := utils.FunctionExists(ctx, a.lambdaClient, functionName)
	if err != nil {
		return nil, fmt.Errorf("checking if function exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("lambda function %q does not exist", functionName)
	}

	exists, err = utils.QualifierExists(ctx, a.lambdaClient, functionName, version)
	if err != nil {
		return nil, fmt.Errorf("checking if version exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("version %q does not exist", version)
	}

	return metrics.GetFunctionConfiguration(ctx, a.lambdaClient, query)
}
