package examples

import (
	"context"
	"time"
	"fmt"

	"github.com/dominikhei/aws-lambda-analyzer/sdk/analyzer"
	sdktypes "github.com/dominikhei/aws-lambda-analyzer/sdk/types"
)

func examples() {
	ctx := context.Background()
    opts := sdktypes.ConfigOptions{
        Region:  "us-west-2",
        Profile: "default",
	}

	a := analyzer.New(ctx, opts)
    functionName := "my-function"
    qualifier := "prod"
    startTime := time.Now().Add(-1 * time.Hour)
    endTime := time.Now()
    period := int32(60)

	result, err := a.GetThrottleRate(ctx, functionName, qualifier, startTime, endTime, period)
		if err != nil {
			fmt.Printf("failed to get throttle rate: %v", err)
		}

    fmt.Printf("Throttle rate for %s:%s between %s and %s is %.2f%%\n",
        result.FunctionName,
        result.Qualifier,
        result.StartTime.Format(time.RFC3339),
        result.EndTime.Format(time.RFC3339),
        result.ThrottleRate*100,
    )
}