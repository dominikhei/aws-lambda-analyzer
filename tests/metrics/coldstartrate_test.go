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

func TestGetColdStartRate_HappyPath(t *testing.T) {
	cw := &mockCWFetcher{
		results: []types.MetricDataResult{
			{Values: []float64{100}},
		},
	}
	logs := &mockLogsFetcher{
		results: []map[string]string{
			{"totalInvocations": "100", "coldStartLines": "10"},
		},
	}

	query := sdktypes.FunctionQuery{
		FunctionName: "test-fn",
		Region:       "us-east-1",
		Qualifier:    "$LATEST",
		StartTime:    time.Now().Add(-1 * time.Hour),
		EndTime:      time.Now(),
	}
	result, err := metrics.GetColdStartRate(context.Background(), logs, cw, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ColdStartRate != 0.1 {
		t.Errorf("expected 0.1, got %v", result.ColdStartRate)
	}
}

func TestGetColdStartRate_NoInvocations(t *testing.T) {
	cw := &mockCWFetcher{
		results: []types.MetricDataResult{
			{Values: []float64{0}},
		},
	}
	logs := &mockLogsFetcher{}

	query := sdktypes.FunctionQuery{
		FunctionName: "test-fn",
		Region:       "us-east-1",
		Qualifier:    "$LATEST",
		StartTime:    time.Now().Add(-1 * time.Hour),
		EndTime:      time.Now(),
	}
	_, err := metrics.GetColdStartRate(context.Background(), logs, cw, query)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var invErr *sdkerrors.NoInvocationsError
	if !errors.As(err, &invErr) {
		t.Errorf("expected NoInvocationsError, got %T", err)
	}
}

// This case is not possible with the AWS API but was added as a caution measure.
func TestGetColdStartRate_EmptyLogData(t *testing.T) {
	cw := &mockCWFetcher{
		results: []types.MetricDataResult{
			{Values: []float64{100}},
		},
	}
	logs := &mockLogsFetcher{
		results: []map[string]string{
			{"totalInvocations": "", "coldStartLines": ""},
		},
	}

	query := sdktypes.FunctionQuery{
		FunctionName: "test-fn",
		Region:       "us-east-1",
		Qualifier:    "$LATEST",
		StartTime:    time.Now().Add(-1 * time.Hour),
		EndTime:      time.Now(),
	}
	result, err := metrics.GetColdStartRate(context.Background(), logs, cw, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ColdStartRate != 0.0 {
		t.Errorf("expected 0.0, got %v", result.ColdStartRate)
	}
}
