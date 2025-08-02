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

package types

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

// ConfigOptions can be used to configure connections to AWS, if the default credentials chain shall be adjusted.
// This can be done by overwriting the default region or using a specific profile or even credentials.
type ConfigOptions struct {
	Region          string
	Profile         string
	AccessKeyID     string
	SecretAccessKey string
}

// FunctionQuery defines the parameters to query metrics for a specific AWS Lambda function.
type FunctionQuery struct {
	FunctionName string    // The name of the Lambda function, e.g., "my-function"
	Region       string    // AWS region, e.g., "us-east-1"
	Qualifier    string    // Lambda version, e.g., "$LATEST", "1"
	StartTime    time.Time // Start of the query interval (UTC)
	EndTime      time.Time // End of the query interval (UTC)
}

// AWSClients holds the clients that are used internally to request AWS Services.
type AWSClients struct {
	LambdaClient     *lambda.Client
	CloudWatchClient *cloudwatch.Client
	LogsClient       *cloudwatchlogs.Client
}

// ThrottleRateReturn is the return of GetThrottleRate.
type ThrottleRateReturn struct {
	ThrottleRate float64   `json:"throttleRate"`
	FunctionName string    `json:"functionName"`
	Qualifier    string    `json:"qualifier"`
	StartTime    time.Time `json:"startTime"`
	EndTime      time.Time `json:"endTime"`
}

// TimeoutRateReturn is the return of GetTimeoutRate.
type TimeoutRateReturn struct {
	TimeoutRate  float64   `json:"timeoutRate"`
	FunctionName string    `json:"functionName"`
	Qualifier    string    `json:"qualifier"`
	StartTime    time.Time `json:"startTime"`
	EndTime      time.Time `json:"endTime"`
}

// ColdStartRateReturn is the return of GetColdStartRate.
type ColdStartRateReturn struct {
	ColdStartRate float64   `json:"coldStartRate"`
	FunctionName  string    `json:"functionName"`
	Qualifier     string    `json:"qualifier"`
	StartTime     time.Time `json:"startTime"`
	EndTime       time.Time `json:"endTime"`
}

// MemoryUsagePercentilesReturn holds various statistics on the maximum used memory of invocations.
// P95UsageRate, P99UsageRate and Conf95UsageRate can be nil if not enough values are present in
// the specified inteval, to calculate them robustly.
type MemoryUsagePercentilesReturn struct {
	MinUsageRate    float64   `json:"minUsageRate"`              // Min (max) Memory usage of any run
	MaxUsageRate    float64   `json:"maxUsageRate"`              // Max (max) Memory usage of any run
	MedianUsageRate float64   `json:"medianUsageRate"`           // Median (max) Memory usage of any run
	MeanUsageRate   float64   `json:"meanUsageRate"`             // Mean (max) Memory usage of any run
	P95UsageRate    *float64  `json:"p95UsageRate,omitempty"`    // 95th percentile
	P99UsageRate    *float64  `json:"p99UsageRate,omitempty"`    // 99th percentile
	Conf95UsageRate *float64  `json:"conf95UsageRate,omitempty"` // 95% confidence interval
	FunctionName    string    `json:"functionName"`
	Qualifier       string    `json:"qualifier"`
	StartTime       time.Time `json:"startTime"`
	EndTime         time.Time `json:"endTime"`
}

// BaseStatisticsReturn contains general statistics on a lambda function.
type BaseStatisticsReturn struct {
	FunctionARN          string            `json:"functionArn"`
	FunctionName         string            `json:"functionName"`
	Qualifier            string            `json:"qualifier"`
	MemorySizeMB         *int32            `json:"memorySizeMb,omitempty"`
	TimeoutSeconds       *int32            `json:"timeoutSeconds,omitempty"`
	Runtime              string            `json:"runtime"`
	LastModified         string            `json:"lastModified"`
	EnvironmentVariables map[string]string `json:"environmentVariables"`
}

// ErrorRateReturn is the return of GetErrorRate.
type ErrorRateReturn struct {
	FunctionName string    `json:"functionName"`
	Qualifier    string    `json:"qualifier"`
	StartTime    time.Time `json:"startTime"`
	EndTime      time.Time `json:"endTime"`
	ErrorRate    float64   `json:"errorRate"`
}

// ErrorType represents a categorized error encountered by an AWS Lambda function.
type ErrorType struct {
	ErrorCategory string `json:"errorCategory"` // ErrorCategory is a semantic extraction what follows after [ERROR] in a log.
	ErrorCount    int    `json:"errorCount"`
}

// ErrorTypesReturn is the return of GetErrorCategoryStatistics.
// It contains a slice of ErrorType.
type ErrorTypesReturn struct {
	Errors       []ErrorType `json:"errors"`
	FunctionName string      `json:"functionName"`
	Qualifier    string      `json:"qualifier"`
	StartTime    time.Time   `json:"startTime"`
	EndTime      time.Time   `json:"endTime"`
}

// DurationStatisticsReturn holds various statistics on the duration of invocations.
// P95Duration, P99Duration and Conf95Duration can be nil if not enough values are present in
// the specified inteval, to calculate them robustly.
type DurationStatisticsReturn struct {
	MinDuration    float64   `json:"minDuration"`              // Min duration of any run
	MaxDuration    float64   `json:"maxDuration"`              // Max duration of any run
	MedianDuration float64   `json:"medianDuration"`           // Median duration of any run
	MeanDuration   float64   `json:"meanDuration"`             // Mean duration of any run
	P95Duration    *float64  `json:"p95Duration,omitempty"`    // 95th percentile duration
	P99Duration    *float64  `json:"p99Duration,omitempty"`    // 99th percentile duration
	Conf95Duration *float64  `json:"conf95Duration,omitempty"` // 95% confidence interval of the durations
	FunctionName   string    `json:"functionName"`
	Qualifier      string    `json:"qualifier"`
	StartTime      time.Time `json:"startTime"`
	EndTime        time.Time `json:"endTime"`
}

// ColdStartDurationStatisticsReturn holds various statistics on the coldstart duration of invocations.
// P95ColdStartDuration, P99ColdStartDuration and Conf95ColdStartDuration can be nil
// if not enough values are present in the specified inteval, to calculate them robustly.
type ColdStartDurationStatisticsReturn struct {
	MinColdStartDuration    float64   `json:"minColdStartDuration"`              // Min coldstart duration of any run
	MaxColdStartDuration    float64   `json:"maxColdStartDuration"`              // Max coldstart duration of any run
	MedianColdStartDuration float64   `json:"medianColdStartDuration"`           // Median coldstart duration of any run
	MeanColdStartDuration   float64   `json:"meanColdStartDuration"`             // Mean coldstart duration of any run
	P95ColdStartDuration    *float64  `json:"p95ColdStartDuration,omitempty"`    // 95th percentile coldstart duration
	P99ColdStartDuration    *float64  `json:"p99ColdStartDuration,omitempty"`    // 99th percentile coldstart duration
	Conf95ColdStartDuration *float64  `json:"conf95ColdStartDuration,omitempty"` // 95% confidence interval of the coldstart durations
	FunctionName            string    `json:"functionName"`
	Qualifier               string    `json:"qualifier"`
	StartTime               time.Time `json:"startTime"`
	EndTime                 time.Time `json:"endTime"`
}

// WasteRatioReturn is the return of GetWasteRatio.
type WasteRatioReturn struct {
	WasteRatio   float64   `json:"wasteRatio"`
	FunctionName string    `json:"functionName"`
	Qualifier    string    `json:"qualifier"`
	StartTime    time.Time `json:"startTime"`
	EndTime      time.Time `json:"endTime"`
}

// Prometheusconfig is used to configure
type PrometheusConfig struct {
	URL      string            `json:"url"`
	JobName  string            `json:"jobName"`
	Grouping map[string]string `json:"grouping"`
	Enabled  bool              `json:"enabled"`
}

// Prometheusconfig is used to configure
type PrometheusConfig struct {
	URL      string
	JobName  string
	Grouping map[string]string
	Enabled  bool
}
