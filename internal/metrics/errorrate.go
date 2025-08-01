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

package metrics

import (
	"context"
	"fmt"

	sdkerrors "github.com/dominikhei/serverless-statistics/errors"
	"github.com/dominikhei/serverless-statistics/internal/cache"
	sdkinterfaces "github.com/dominikhei/serverless-statistics/internal/interfaces"
	sdktypes "github.com/dominikhei/serverless-statistics/types"

	"github.com/dominikhei/serverless-statistics/internal/utils"
)

// GetErrorRate calculates the ratio of errors and total invoations
// of an AWS Lambda function over a specified time range and qualifier (version).
func GetErrorRate(
	ctx context.Context,
	cwFetcher sdkinterfaces.CloudWatchFetcher,
	invocationsCache sdkinterfaces.Cache,
	query sdktypes.FunctionQuery,
) (*sdktypes.ErrorRateReturn, error) {

	// cache reduces the number of calls to CloudWatch metrics.
	// It lives as long as the Go process is running.
	key := cache.CacheKey{
		FunctionName: query.FunctionName,
		Qualifier:    query.Qualifier,
		Start:        query.StartTime,
		End:          query.EndTime,
	}
	var invocationsSum float64
	if invocationsCache.Has(key) {
		invocations, _ := invocationsCache.Get(key)
		invocationsSum = float64(invocations)
	} else {
		invocationsResults, err := cwFetcher.FetchMetric(ctx, query, "Invocations", "Sum")
		if err != nil {
			return nil, fmt.Errorf("fetch invocations metric: %w", err)
		}
		invocationsSum, err = utils.SumMetricValues(invocationsResults)
		if err != nil {
			return nil, fmt.Errorf("parse invocations metric data: %w", err)
		}
		invocationsCache.Set(key, int(invocationsSum))
	}
	if invocationsSum == 0 {
		return nil, &sdkerrors.NoInvocationsError{FunctionName: query.FunctionName}
	}

	errorsResults, err := cwFetcher.FetchMetric(ctx, query, "Errors", "Sum")
	if err != nil {
		return nil, fmt.Errorf("fetch errors metric: %w", err)
	}
	errorsSum, err := utils.SumMetricValues(errorsResults)
	if err != nil {
		return nil, fmt.Errorf("parse errors metric data: %w", err)
	}

	errorRate := float64(errorsSum) / float64(invocationsSum)

	return &sdktypes.ErrorRateReturn{
		FunctionName: query.FunctionName,
		Qualifier:    query.Qualifier,
		StartTime:    query.StartTime,
		EndTime:      query.EndTime,
		ErrorRate:    errorRate,
	}, nil
}
