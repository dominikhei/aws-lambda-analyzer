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

	cwTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	sdkerrors "github.com/dominikhei/serverless-statistics/errors"
	"github.com/dominikhei/serverless-statistics/internal/cache"
	"github.com/dominikhei/serverless-statistics/internal/metrics"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

func TestGetThrottleRate_HappyPath(t *testing.T) {
	mock := &mockCWFetcher{
		results: []cwTypes.MetricDataResult{
			{Values: []float64{50}}, // both Invocations and Throttles
		},
	}

	cache := cache.NewCache()

	query := sdktypes.FunctionQuery{
		FunctionName: "test-fn",
		Region:       "us-east-1",
		Qualifier:    "1",
		StartTime:    time.Now().Add(-10 * time.Minute),
		EndTime:      time.Now(),
	}

	// First call: cache miss, will fetch metrics
	result, err := metrics.GetThrottleRate(context.Background(), mock, cache, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedRate := 1.0
	if result.ThrottleRate != expectedRate {
		t.Errorf("expected throttle rate %v, got %v", expectedRate, result.ThrottleRate)
	}

	result, err = metrics.GetThrottleRate(context.Background(), mock, cache, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ThrottleRate != expectedRate {
		t.Errorf("expected throttle rate %v, got %v", expectedRate, result.ThrottleRate)
	}
}

func TestGetThrottleRate_NoInvocations(t *testing.T) {
	mock := &mockCWFetcher{
		results: []cwTypes.MetricDataResult{
			{Values: []float64{0}}, // zero invocations
		},
		err: nil,
	}

	cache := cache.NewCache()

	query := sdktypes.FunctionQuery{
		FunctionName: "test-fn",
		Region:       "us-east-1",
	}

	_, err := metrics.GetThrottleRate(context.Background(), mock, cache, query)
	if err == nil {
		t.Fatal("expected error for no invocations, got nil")
	}
	var noInvErr *sdkerrors.NoInvocationsError
	if !errors.As(err, &noInvErr) {
		t.Errorf("expected NoInvocationsError, got: %v", err)
	}
}
