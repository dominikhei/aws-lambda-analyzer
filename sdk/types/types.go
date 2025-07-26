package types 

import  (
	"time"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
    "github.com/aws/aws-sdk-go-v2/service/xray"
)
// FunctionQuery defines the parameters to query metrics for a specific AWS Lambda function.
type FunctionQuery struct {
    FunctionName string    // The name of the Lambda function, e.g., "my-function"
    Region       string    // AWS region, e.g., "us-east-1"
    Qualifier    string    // Lambda version or alias, e.g., "$LATEST", "1", or "prod"
    StartTime    time.Time // Start of the query interval (UTC)
    EndTime      time.Time // End of the query interval (UTC)
}

type AWSClients struct {
	LambdaClient     *lambda.Client
	CloudWatchClient *cloudwatch.Client
	XRayClient       *xray.Client
	LogsClient       *cloudwatchlogs.Client
}

type ThrottleRateReturn struct {
	ThrottleRate float64
	FunctionName string
	Qualifier string
    StartTime    time.Time // Start of the query interval (UTC)
    EndTime      time.Time // End of the query interval (UTC)	
}

type TimeoutRateReturn struct {
	TimeoutRate float64
	FunctionName string
	Qualifier string
    StartTime    time.Time // Start of the query interval (UTC)
    EndTime      time.Time // End of the query interval (UTC)	
}

type ConfigOptions struct {
	Region          string
	Profile         string
	AccessKeyID     string
	SecretAccessKey string
}

type ColdStartRateReturn struct {
	ColdStartRate float32 // Timedout Invocations / Total
	FunctionName string // Name of the function
	Qualifier string // Qualifier of the function
    StartTime    time.Time // Start of the query interval (UTC)
    EndTime      time.Time // End of the query interval (UTC)	
}

// MemoryUsagePercentilesReturn holds various statistics on the maximum used memory of invocations
type MemoryUsagePercentilesReturn struct{
        MinUsageRate    float32 // Min (max) Memory usage of any run
        MaxUsageRate    float32 // Max (max) Memory usage of any run
        MedianUsageRate float32 // Median (max) Memory usage of any run
		MeanUsageRate   float32 // Mean (max) Memory usage of any run
        P95UsageRate    *float32 // Pointers as these values can be nil
        P99UsageRate    *float32 // in case of too little samples
		Conf95UsageRate *float32 // 95% confidence interval
        FunctionARN    string // ARN of the lambda function
        Qualifier       string // Qualifier of the lambda function
        StartTime       time.Time // earliest considered invocation
        EndTime         time.Time // latest considered invocation
    }