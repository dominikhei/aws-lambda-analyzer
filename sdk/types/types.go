package types 

import "time"

// FunctionQuery defines the parameters to query metrics for a specific AWS Lambda function.
type FunctionQuery struct {
    FunctionName string    // The name of the Lambda function, e.g., "my-function"
    Region       string    // AWS region, e.g., "us-east-1"
    Qualifier    string    // Lambda version or alias, e.g., "$LATEST", "1", or "prod"
    StartTime    time.Time // Start of the query interval (UTC)
    EndTime      time.Time // End of the query interval (UTC)
}