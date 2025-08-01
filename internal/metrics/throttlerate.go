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
	"github.com/dominikhei/serverless-statistics/internal/utils"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

// GetThrottleRate calculates the throttle rate of an AWS Lambda function
// over a specified time range and qualifier (version).
// The throttle rate is computed as throttled invocations divided by total invocations.
func GetThrottleRate(
	ctx context.Context,
	cwFetcher sdkinterfaces.CloudWatchFetcher,
	invocationsCache sdkinterfaces.Cache,
	query sdktypes.FunctionQuery,
) (*sdktypes.ThrottleRateReturn, error) {

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

	throttlesResults, err := cwFetcher.FetchMetric(ctx, query, "Throttles", "Sum")
	if err != nil {
		return nil, fmt.Errorf("fetch throttles metric: %w", err)
	}
	throttlesSum, err := utils.SumMetricValues(throttlesResults)
	if err != nil {
		return nil, fmt.Errorf("parse throttles metric data: %w", err)
	}

	throttleRate := throttlesSum / invocationsSum
	result := &sdktypes.ThrottleRateReturn{
		ThrottleRate: throttleRate,
		FunctionName: query.FunctionName,
		Qualifier:    query.Qualifier,
		StartTime:    query.StartTime,
		EndTime:      query.EndTime,
	}
	return result, nil
}
