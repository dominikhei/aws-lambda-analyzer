package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	sdkerrors "github.com/dominikhei/serverless-statistics/errors"
	"github.com/dominikhei/serverless-statistics/internal/metrics"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

func TestGetErrorRate_HappyPath(t *testing.T) {
	mockCW := &mockCWFetcher{
		results: []types.MetricDataResult{
			{Values: []float64{50}},
		},
		err: nil,
	}
	query := sdktypes.FunctionQuery{
		FunctionName: "test-fn",
		Region:       "us-east-1",
		Qualifier:    "1",
		StartTime:    time.Now().Add(-10 * time.Minute),
		EndTime:      time.Now(),
	}

	result, err := metrics.GetErrorRate(context.Background(), mockCW, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ErrorRate != 1.0 {
		t.Errorf("expected error rate 1.0, got %v", result.ErrorRate)
	}
}

func TestGetErrorRate_NoInvocations(t *testing.T) {
	mockCW := &mockCWFetcher{
		results: []types.MetricDataResult{
			{Values: []float64{0}},
		},
		err: nil,
	}
	query := sdktypes.FunctionQuery{
		FunctionName: "test-fn",
		Region:       "us-east-1",
		Qualifier:    "1",
		StartTime:    time.Now().Add(-10 * time.Minute),
		EndTime:      time.Now(),
	}

	_, err := metrics.GetErrorRate(context.Background(), mockCW, query)
	if err == nil {
		t.Fatal("expected NoInvocationsError, got nil")
	}
	var noInvErr *sdkerrors.NoInvocationsError
	if !errors.As(err, &noInvErr) {
		t.Errorf("expected NoInvocationsError, got: %v", err)
	}
}
