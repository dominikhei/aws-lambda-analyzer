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
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	sdkerrors "github.com/dominikhei/serverless-statistics/errors"
	"github.com/dominikhei/serverless-statistics/internal/metrics"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

func TestGetTimeoutRate_HappyPath(t *testing.T) {
	cw := &mockCWFetcher{
		results: []types.MetricDataResult{
			{Values: []float64{100}},
		},
		err: nil,
	}
	logs := &mockLogsFetcher{
		results: []map[string]string{
			{"invocationsCount": "100", "timeoutCount": "10"},
		},
		err: nil,
	}

	query := sdktypes.FunctionQuery{
		FunctionName: "test-fn",
		Region:       "us-east-1",
		Qualifier:    "$LATEST",
		StartTime:    time.Now().Add(-15 * time.Minute),
		EndTime:      time.Now(),
	}

	result, err := metrics.GetTimeoutRate(context.Background(), cw, logs, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := 10.0 / 100.0
	if result.TimeoutRate != expected {
		t.Errorf("expected timeout rate %v, got %v", expected, result.TimeoutRate)
	}
	if result.FunctionName != query.FunctionName {
		t.Errorf("expected function name %s, got %s", query.FunctionName, result.FunctionName)
	}
}

func TestGetTimeoutRate_NoInvocations(t *testing.T) {
	cw := &mockCWFetcher{
		results: []types.MetricDataResult{
			{Values: []float64{0}},
		},
		err: nil,
	}
	query := sdktypes.FunctionQuery{
		FunctionName: "test-fn",
	}
	logs := &mockLogsFetcher{}

	_, err := metrics.GetTimeoutRate(context.Background(), cw, logs, query)
	if err == nil {
		t.Fatal("expected error for no invocations, got nil")
	}
	var noInvErr *sdkerrors.NoInvocationsError
	if !errors.As(err, &noInvErr) {
		t.Errorf("expected NoInvocationsError, got: %v", err)
	}
}

func TestGetTimeoutRate_ParseInvocationCountError(t *testing.T) {
	cw := &mockCWFetcher{
		results: []types.MetricDataResult{{Values: []float64{10}}},
	}
	logs := &mockLogsFetcher{
		results: []map[string]string{
			{"invocationsCount": "NaN", "timeoutCount": "10"},
		},
		err: nil,
	}

	query := sdktypes.FunctionQuery{
		FunctionName: "test-fn",
		Region:       "us-east-1",
		Qualifier:    "$LATEST",
		StartTime:    time.Now().Add(-15 * time.Minute),
		EndTime:      time.Now(),
	}

	_, err := metrics.GetTimeoutRate(context.Background(), cw, logs, query)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "parse invocationsCount from logs") {
		t.Fatalf("expected error containing 'parse invocationsCount from logs', got: %v", err)
	}
}

func TestGetTimeoutRate_ZeroTimeoutCount(t *testing.T) {
	cw := &mockCWFetcher{
		results: []types.MetricDataResult{{Values: []float64{10}}},
	}
	logs := &mockLogsFetcher{
		results: []map[string]string{
			{"invocationsCount": "10", "timeoutCount": "0"},
		},
		err: nil,
	}

	query := sdktypes.FunctionQuery{
		FunctionName: "test-fn",
		Region:       "us-east-1",
		Qualifier:    "$LATEST",
		StartTime:    time.Now().Add(-15 * time.Minute),
		EndTime:      time.Now(),
	}

	result, err := metrics.GetTimeoutRate(context.Background(), cw, logs, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TimeoutRate != 0.0 {
		t.Errorf("expected timeout rate 0.0, got %v", result.TimeoutRate)
	}
}
