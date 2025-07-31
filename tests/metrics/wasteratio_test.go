// Copyright 2025 dominikhei
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

func TestGetWasteRatio_HappyPath(t *testing.T) {
	cw := &mockCWFetcher{
		results: []types.MetricDataResult{
			{Values: []float64{100}},
		},
	}
	logs := &mockLogsFetcher{
		results: []map[string]string{
			{"totalDuration": "100", "totalBilledDuration": "110"},
		},
	}

	query := sdktypes.FunctionQuery{
		FunctionName: "test-fn",
		Region:       "us-east-1",
		Qualifier:    "$LATEST",
		StartTime:    time.Now().Add(-1 * time.Hour),
		EndTime:      time.Now(),
	}
	result, err := metrics.GetWasteRatio(context.Background(), cw, logs, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.WasteRatio != 0.09090909090909091 {
		t.Errorf("expected 0.09090909090909091, got %v", result.WasteRatio)
	}
}

func TestGetWasteRatio_NoInvocations(t *testing.T) {
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
	_, err := metrics.GetWasteRatio(context.Background(), cw, logs, query)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var invErr *sdkerrors.NoInvocationsError
	if !errors.As(err, &invErr) {
		t.Errorf("expected NoInvocationsError, got %T", err)
	}
}

// This case is not possible with the AWS API but was added as a caution measure.
func TestGetWasteRatio_EmptyLogData(t *testing.T) {
	cw := &mockCWFetcher{
		results: []types.MetricDataResult{
			{Values: []float64{100}},
		},
	}
	logs := &mockLogsFetcher{
		results: []map[string]string{
			{"totalDuration": "", "totalBilledDuration": ""},
		},
	}

	query := sdktypes.FunctionQuery{
		FunctionName: "test-fn",
		Region:       "us-east-1",
		Qualifier:    "$LATEST",
		StartTime:    time.Now().Add(-1 * time.Hour),
		EndTime:      time.Now(),
	}
	_, err := metrics.GetWasteRatio(context.Background(), cw, logs, query)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	expected := "total duration is zero, cannot calculate waste ratio"
	if err.Error() != expected {
		t.Errorf("unexpected error: got %q, want %q", err.Error(), expected)
	}
}
