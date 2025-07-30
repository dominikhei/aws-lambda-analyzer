package main

import (
	"context"
	"fmt"
	"time"

	serverlessstatistics "github.com/dominikhei/serverless-statistics"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

func main() {
	ctx := context.Background()
	opts := sdktypes.ConfigOptions{
		Region:  "eu-central-1",
		Profile: "default",
	}

	a := serverlessstatistics.New(ctx, opts)
	functionName := "bitvavo"
	qualifier := ""
	layout := "2006-01-02 15:04:05"
	startTime, _ := time.Parse(layout, "2025-07-17 01:00:00")
	endTime := time.Now()
	rate, err := a.GetTimeoutRate(ctx, functionName, qualifier, startTime, endTime)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Print(rate)
}
