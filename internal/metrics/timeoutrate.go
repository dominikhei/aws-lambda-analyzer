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
	"strconv"
	"strings"

	sdkerrors "github.com/dominikhei/serverless-statistics/errors"
	"github.com/dominikhei/serverless-statistics/internal/cache"
	sdkinterfaces "github.com/dominikhei/serverless-statistics/internal/interfaces"
	"github.com/dominikhei/serverless-statistics/internal/queries"
	"github.com/dominikhei/serverless-statistics/internal/utils"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

// GetTimeoutRate calculates the timeout rate of an AWS Lambda function
// over a specified time range and qualifier (version).
// The timeout rate is computed as timed-out invocations divided by total invocations.
func GetTimeoutRate(
	ctx context.Context,
	cwFetcher sdkinterfaces.CloudWatchFetcher,
	logsFetcher sdkinterfaces.LogsInsightsFetcher,
	invocationsCache sdkinterfaces.Cache,
	query sdktypes.FunctionQuery,
) (*sdktypes.TimeoutRateReturn, error) {

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

	escapedQualifier := strings.ReplaceAll(query.Qualifier, "$", "\\$")
	queryString := fmt.Sprintf(queries.LambdaUniqueRequestsWithVersion, escapedQualifier)
	results, err := logsFetcher.RunQuery(ctx, query, queryString)
	if err != nil {
		return nil, fmt.Errorf("fetch errors from logs insights: %w", err)
	}
	var invocationsCount int
	if len(results) > 0 {
		val := results[0]["invocationsCount"]
		if val != "" {
			invocationsCount, err = strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("parse invocationsCount from logs: %w", err)
			}
		}
	}
	escapedQualifier = strings.ReplaceAll(query.Qualifier, "$", "\\$")
	queryString = fmt.Sprintf(queries.LambdaTimeoutQueryWithVersion, escapedQualifier)
	results, err = logsFetcher.RunQuery(ctx, query, queryString)
	if err != nil {
		return nil, fmt.Errorf("run logs insights query: %w", err)
	}
	var timeoutCount int
	if len(results) > 0 {
		val := results[0]["timeoutCount"]
		if val != "" {
			timeoutCount, err = strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("parse timeoutCount from logs: %w", err)
			}
		}
	}
	timeoutRate := float64(timeoutCount) / float64(invocationsCount)

	return &sdktypes.TimeoutRateReturn{
		TimeoutRate:  timeoutRate,
		FunctionName: query.FunctionName,
		Qualifier:    query.Qualifier,
		StartTime:    query.StartTime,
		EndTime:      query.EndTime,
	}, nil
}
