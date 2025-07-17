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