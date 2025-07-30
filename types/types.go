package types

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/xray"
)

// Configoptions can be used to configure connections to AWS, if the default credentials chain shall be adjusted.
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
	XRayClient       *xray.Client
	LogsClient       *cloudwatchlogs.Client
}

// ThrottleRateReturn is the return of GetThrottleRate.
type ThrottleRateReturn struct {
	ThrottleRate float64
	FunctionName string
	Qualifier    string
	StartTime    time.Time
	EndTime      time.Time
}

// TimeoutRateReturn is the return of GetTimeoutRate.
type TimeoutRateReturn struct {
	TimeoutRate  float64
	FunctionName string
	Qualifier    string
	StartTime    time.Time
	EndTime      time.Time
}

// ColdStartRateReturn is the return of GetColdStartRate.
type ColdStartRateReturn struct {
	ColdStartRate float32   // Timedout Invocations / Total
	FunctionName  string    // Name of the function
	Qualifier     string    // Qualifier of the function
	StartTime     time.Time // Start of the query interval (UTC)
	EndTime       time.Time // End of the query interval (UTC)
}

// MemoryUsagePercentilesReturn holds various statistics on the maximum used memory of invocations.
// P95UsageRate, P99UsageRate and Conf95UsageRate can be nil if not enough values are present in
// the specified inteval, to calculate them robustly.
type MemoryUsagePercentilesReturn struct {
	MinUsageRate    float32  // Min (max) Memory usage of any run
	MaxUsageRate    float32  // Max (max) Memory usage of any run
	MedianUsageRate float32  // Median (max) Memory usage of any run
	MeanUsageRate   float32  // Mean (max) Memory usage of any run
	P95UsageRate    *float32 // 95th percentile
	P99UsageRate    *float32 // 99th percentile
	Conf95UsageRate *float32 // 95% confidence interval
	FunctionName    string
	Qualifier       string
	StartTime       time.Time
	EndTime         time.Time
}

// BaseStatistics contains general statistics on the lambda function.
type BaseStatisticsReturn struct {
	FunctionARN            string
	FunctionName           string
	Qualifier              string
	MemorySizeMB           int
	TimeoutSeconds         int
	Runtime                string
	LastModified           string
	ProvisionedConcurrency int
	NumInvocations         int64
	EnvironmentVariables   map[string]string
	StartTime              time.Time
	EndTime                time.Time
}

// ErrorRateReturn is the return of GetErrorRate.
type ErrorRateReturn struct {
	FunctionName string
	Qualifier    string
	StartTime    time.Time
	EndTime      time.Time
	ErrorRate    float32
}

// ErrorType represents a categorized error encountered by an AWS Lambda function.
type ErrorType struct {
	ErrorCategory string // ErrorCategory is a semantic extraction what follows after [ERROR] in a log.
	ErrorCount    int64
}

// ErrorTypesReturn is the return of GetErrorCategoryStatistics.
// It contains a slice of ErrorType.
type ErrorTypesReturn struct {
	Errors       []ErrorType
	FunctionName string
	Qualifier    string
	StartTime    time.Time
	EndTime      time.Time
}

// MemoryUsagePercentilesReturn holds various statistics on the maximum used memory of invocations.
// P95UsageRate, P99UsageRate and Conf95UsageRate can be nil if not enough values are present in
// the specified inteval, to calculate them robustly.
type DurationStatisticsReturn struct {
	MinDuration    float32  // Min duration of any run
	MaxDuration    float32  // Max duration of any run
	MedianDuration float32  // Median duration of any run
	MeanDuration   float32  // Mean duration of any run
	P95Duration    *float32 // 95th percentile duration
	P99Duration    *float32 // 99th percentile duration
	Conf95Duration *float32 // 95% confidence interval of the durations
	FunctionName   string
	Qualifier      string
	StartTime      time.Time
	EndTime        time.Time
}
