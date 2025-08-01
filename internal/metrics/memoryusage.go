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

// GetMaxMemoryUsageStatistics calculates statistics on the maximum used memory
// of an AWS Lambda functions invocations over a specified time range and qualifier (version).
func GetMaxMemoryUsageStatistics(
	ctx context.Context,
	logsFetcher sdkinterfaces.LogsInsightsFetcher,
	cwFetcher sdkinterfaces.CloudWatchFetcher,
	invocationsCache sdkinterfaces.Cache,
	query sdktypes.FunctionQuery,
) (*sdktypes.MemoryUsagePercentilesReturn, error) {

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
	queryString := fmt.Sprintf(queries.LambdaMemoryUtilizationQueryWithVersion, escapedQualifier)
	results, err := logsFetcher.RunQuery(ctx, query, queryString)
	if err != nil {
		return nil, fmt.Errorf("run logs insights query: %w", err)
	}

	var ratios []float64
	for _, row := range results {
		if valStr, ok := row["memoryUtilizationRatio"]; ok {
			if val, err := strconv.ParseFloat(valStr, 64); err == nil {
				ratios = append(ratios, val)
			} else {
				fmt.Printf("warn: could not parse %q as float64: %v", valStr, err)
			}
		}
	}
	memoryStats, err := utils.CalcSummaryStats(ratios)
	if err != nil {
		return nil, fmt.Errorf("error calculating summary statistics: %w", err)
	}
	return &sdktypes.MemoryUsagePercentilesReturn{
		MinUsageRate:    memoryStats.Min,
		MaxUsageRate:    memoryStats.Max,
		MedianUsageRate: memoryStats.Median,
		MeanUsageRate:   memoryStats.Mean,
		P95UsageRate:    memoryStats.P95,
		P99UsageRate:    memoryStats.P99,
		Conf95UsageRate: memoryStats.ConfInt95,
		FunctionName:    query.FunctionName,
		Qualifier:       query.Qualifier,
		StartTime:       query.StartTime,
		EndTime:         query.EndTime,
	}, nil
}
