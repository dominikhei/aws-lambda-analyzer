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

func TestGetErrorTypes_HappyPath(t *testing.T) {
	cw := &mockCWFetcher{
		results: []types.MetricDataResult{{Values: []float64{10}}},
	}
	logs := &mockLogsFetcher{
		results: []map[string]string{
			{"error_category": "TimeoutError", "error_count": "5"},
			{"error_category": "ValidationError", "error_count": "3"},
		},
	}

	query := sdktypes.FunctionQuery{
		FunctionName: "my-function",
		Qualifier:    "$LATEST",
		StartTime:    time.Now().Add(-1 * time.Hour),
		EndTime:      time.Now(),
	}

	result, err := metrics.GetErrorTypes(context.Background(), logs, cw, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Errors) != 2 {
		t.Fatalf("expected 2 error categories, got %d", len(result.Errors))
	}
	if result.Errors[0].ErrorCategory != "TimeoutError" || result.Errors[0].ErrorCount != 5 {
		t.Errorf("unexpected first error: %+v", result.Errors[0])
	}
	if result.Errors[1].ErrorCategory != "ValidationError" || result.Errors[1].ErrorCount != 3 {
		t.Errorf("unexpected second error: %+v", result.Errors[1])
	}
}

func TestGetErrorTypes_NoInvocations(t *testing.T) {
	cw := &mockCWFetcher{
		results: []types.MetricDataResult{{Values: []float64{0}}},
	}
	logs := &mockLogsFetcher{}

	query := sdktypes.FunctionQuery{
		FunctionName: "empty-fn",
		Qualifier:    "1",
		StartTime:    time.Now().Add(-15 * time.Minute),
		EndTime:      time.Now(),
	}

	_, err := metrics.GetErrorTypes(context.Background(), logs, cw, query)
	if err == nil {
		t.Fatal("expected NoInvocationsError, got nil")
	}
	var noInvErr *sdkerrors.NoInvocationsError
	if !errors.As(err, &noInvErr) {
		t.Errorf("expected NoInvocationsError, got: %v", err)
	}
}

func TestGetErrorTypes_InvalidErrorCount(t *testing.T) {
	cw := &mockCWFetcher{
		results: []types.MetricDataResult{{Values: []float64{5}}},
	}
	logs := &mockLogsFetcher{
		results: []map[string]string{
			{"error_category": "TimeoutError", "error_count": "invalid"},
			{"error_category": "ValidationError", "error_count": "7"},
		},
	}

	query := sdktypes.FunctionQuery{
		FunctionName: "broken-fn",
		Qualifier:    "1",
		StartTime:    time.Now().Add(-1 * time.Hour),
		EndTime:      time.Now(),
	}

	result, err := metrics.GetErrorTypes(context.Background(), logs, cw, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 valid error, got %d", len(result.Errors))
	}
	if result.Errors[0].ErrorCategory != "ValidationError" || result.Errors[0].ErrorCount != 7 {
		t.Errorf("unexpected error: %+v", result.Errors[0])
	}
}

func TestGetErrorTypes_MissingErrorCategory(t *testing.T) {
	cw := &mockCWFetcher{
		results: []types.MetricDataResult{{Values: []float64{5}}},
	}
	logs := &mockLogsFetcher{
		results: []map[string]string{
			{"error_category": "", "error_count": "4"},
			{"error_count": "6"},
		},
	}

	query := sdktypes.FunctionQuery{
		FunctionName: "missing-cat-fn",
		Qualifier:    "1",
		StartTime:    time.Now().Add(-1 * time.Hour),
		EndTime:      time.Now(),
	}

	result, err := metrics.GetErrorTypes(context.Background(), logs, cw, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Errors) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(result.Errors))
	}
	for _, e := range result.Errors {
		if e.ErrorCategory != "UnknownError" {
			t.Errorf("expected UnknownError category for empty/missing, got %q", e.ErrorCategory)
		}
	}
}
