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
	"github.com/dominikhei/serverless-statistics/internal/cache"
	"github.com/dominikhei/serverless-statistics/internal/metrics"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

// Tests for >= 20 invocations to calculate the percentiles will not be added, as this
// logic is already tested in the utils tests.

func TestGetMaxMemoryUsageStatistics_HappyPath(t *testing.T) {
	cw := &mockCWFetcher{
		results: []types.MetricDataResult{
			{Values: []float64{42}},
		},
	}
	logs := &mockLogsFetcher{
		results: []map[string]string{
			{"memoryUtilizationRatio": "0.1"},
			{"memoryUtilizationRatio": "0.5"},
			{"memoryUtilizationRatio": "0.9"},
		},
	}
	cache := cache.NewCache()

	query := sdktypes.FunctionQuery{
		FunctionName: "test-fn",
		Region:       "us-east-1",
		Qualifier:    "$LATEST",
		StartTime:    time.Now().Add(-15 * time.Minute),
		EndTime:      time.Now(),
	}

	result, err := metrics.GetMaxMemoryUsageStatistics(context.Background(), logs, cw, cache, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.MinUsageRate != 0.1 || result.MaxUsageRate != 0.9 || result.MeanUsageRate != 0.5 || result.P95UsageRate != nil || result.P99UsageRate != nil || result.Conf95UsageRate != nil {
		t.Errorf("unexpected stats: %+v", result)
	}
	if result.FunctionName != "test-fn" || result.Qualifier != "$LATEST" {
		t.Errorf("unexpected function metadata: %+v", result)
	}
}

func TestGetMaxMemoryUsageStatistics_NoInvocations(t *testing.T) {
	cw := &mockCWFetcher{
		results: []types.MetricDataResult{
			{Values: []float64{0}},
		},
	}
	logs := &mockLogsFetcher{}
	cache := cache.NewCache()

	query := sdktypes.FunctionQuery{
		FunctionName: "empty-fn",
		Region:       "us-east-1",
		Qualifier:    "1",
		StartTime:    time.Now().Add(-15 * time.Minute),
		EndTime:      time.Now(),
	}

	_, err := metrics.GetMaxMemoryUsageStatistics(context.Background(), logs, cw, cache, query)
	if err == nil {
		t.Fatal("expected NoInvocationsError, got nil")
	}
	var noInvErr *sdkerrors.NoInvocationsError
	if !errors.As(err, &noInvErr) {
		t.Errorf("expected NoInvocationsError, got: %v", err)
	}
}

// This test case is not possible with the AWS API but was added as a caution measure.
func TestGetMaxMemoryUsageStatistics_InvalidMemoryUtilizationEntry(t *testing.T) {
	cw := &mockCWFetcher{
		results: []types.MetricDataResult{
			{Values: []float64{10}},
		},
	}
	logs := &mockLogsFetcher{
		results: []map[string]string{
			{"memoryUtilizationRatio": "invalid"},
			{"memoryUtilizationRatio": "0.7"},
		},
	}
	cache := cache.NewCache()

	query := sdktypes.FunctionQuery{
		FunctionName: "broken-fn",
		Region:       "us-east-1",
		Qualifier:    "1",
		StartTime:    time.Now().Add(-1 * time.Hour),
		EndTime:      time.Now(),
	}

	result, err := metrics.GetMaxMemoryUsageStatistics(context.Background(), logs, cw, cache, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.MinUsageRate != 0.7 || result.MaxUsageRate != 0.7 || result.MeanUsageRate != 0.7 {
		t.Errorf("expected single valid memory usage 0.7, got: %+v", result)
	}
}
