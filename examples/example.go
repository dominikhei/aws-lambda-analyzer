package main

import (
	"context"
	"time"
	"fmt"

	"github.com/dominikhei/aws-lambda-analyzer/sdk/serverlessstatistics"
	sdktypes "github.com/dominikhei/aws-lambda-analyzer/sdk/types"
)

func main() {
	ctx := context.Background()
    opts := sdktypes.ConfigOptions{
        Region:  "eu-central-1",
        Profile: "default",
	}

	a := serverlessstatistics.New(ctx, opts)
    functionName := "bitvavo"
    //qualifier := "prod"
    layout := "2006-01-02 15:04:05"
    startTime, _ := time.Parse(layout, "2025-01-07 00:00:00")
    endTime := time.Now()
    period := int32(1)

	rate, err := a.GetErrorCategoryStatistics(ctx, functionName, "", startTime, endTime, period)
		if err != nil {
			fmt.Printf("failed to get throttle rate: %v", err)
		}
	fmt.Print(rate)

}