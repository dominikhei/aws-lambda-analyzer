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

func TestGetErrorRate_HappyPath(t *testing.T) {
	mockCW := &mockCWFetcher{
		results: []types.MetricDataResult{
			{Values: []float64{50}},
		},
		err: nil,
	}
	cache := cache.NewCache()

	query := sdktypes.FunctionQuery{
		FunctionName: "test-fn",
		Region:       "us-east-1",
		Qualifier:    "1",
		StartTime:    time.Now().Add(-10 * time.Minute),
		EndTime:      time.Now(),
	}

	result, err := metrics.GetErrorRate(context.Background(), mockCW, cache, query)
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
	cache := cache.NewCache()

	query := sdktypes.FunctionQuery{
		FunctionName: "test-fn",
		Region:       "us-east-1",
		Qualifier:    "1",
		StartTime:    time.Now().Add(-10 * time.Minute),
		EndTime:      time.Now(),
	}

	_, err := metrics.GetErrorRate(context.Background(), mockCW, cache, query)
	if err == nil {
		t.Fatal("expected NoInvocationsError, got nil")
	}
	var noInvErr *sdkerrors.NoInvocationsError
	if !errors.As(err, &noInvErr) {
		t.Errorf("expected NoInvocationsError, got: %v", err)
	}
}
