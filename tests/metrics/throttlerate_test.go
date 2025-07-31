package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	cwTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	sdkerrors "github.com/dominikhei/serverless-statistics/errors"
	"github.com/dominikhei/serverless-statistics/internal/metrics"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

func TestGetThrottleRate_HappyPath(t *testing.T) {
	mock := &mockCWFetcher{
		results: []cwTypes.MetricDataResult{
			{Values: []float64{50}}, // This will be returned for both Invocations and Throttles
		},
	}

	query := sdktypes.FunctionQuery{
		FunctionName: "test-fn",
		Region:       "us-east-1",
		Qualifier:    "1",
		StartTime:    time.Now().Add(-10 * time.Minute),
		EndTime:      time.Now(),
	}

	result, err := metrics.GetThrottleRate(context.Background(), mock, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedRate := 1.0 // 50 / 50
	if result.ThrottleRate != expectedRate {
		t.Errorf("expected throttle rate %v, got %v", expectedRate, result.ThrottleRate)
	}
}

func TestGetThrottleRate_NoInvocations(t *testing.T) {
	cw := &mockCWFetcher{
		results: []cwTypes.MetricDataResult{
			{Values: []float64{0}}, // zero invocations
		},
		err: nil,
	}

	query := sdktypes.FunctionQuery{
		FunctionName: "test-fn",
	}

	_, err := metrics.GetThrottleRate(context.Background(), cw, query)
	if err == nil {
		t.Fatal("expected error for no invocations, got nil")
	}
	var noInvErr *sdkerrors.NoInvocationsError
	if !errors.As(err, &noInvErr) {
		t.Errorf("expected NoInvocationsError, got: %v", err)
	}
}
