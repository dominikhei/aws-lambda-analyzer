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
	sdktypes "github.com/dominikhei/serverless-statistics/types"

	"github.com/dominikhei/serverless-statistics/internal/utils"
)

// GetErrorTypes counts the different errors that occur over a specified time range and qualifier (version).
// It uses a regex to search for the error in an [ERROR] line in logs and groups them
// based on semantics.
func GetErrorTypes(
	ctx context.Context,
	logsFetcher sdkinterfaces.LogsInsightsFetcher,
	cwFetcher sdkinterfaces.CloudWatchFetcher,
	invocationsCache sdkinterfaces.Cache,
	query sdktypes.FunctionQuery,
) (*sdktypes.ErrorTypesReturn, error) {

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
	queryString := fmt.Sprintf(queries.LambdaErrorTypesQueryWithVersion, escapedQualifier)
	results, err := logsFetcher.RunQuery(ctx, query, queryString)
	if err != nil {
		return nil, fmt.Errorf("run logs insights query: %w", err)
	}

	var stats []sdktypes.ErrorType
	for _, row := range results {
		category, ok := row["error_category"]
		if !ok || category == "" {
			category = "UnknownError"
		}
		countStr, ok := row["error_count"]
		if !ok {
			continue
		}
		count, err := strconv.ParseInt(countStr, 10, 0)
		if err != nil {
			fmt.Printf("warn: could not parse error_count %q: %v\n", countStr, err)
			continue
		}
		stats = append(stats, sdktypes.ErrorType{
			ErrorCategory: category,
			ErrorCount:    int(count),
		})
	}

	return &sdktypes.ErrorTypesReturn{
		Errors:       stats,
		FunctionName: query.FunctionName,
		Qualifier:    query.Qualifier,
		StartTime:    query.StartTime,
		EndTime:      query.EndTime,
	}, nil
}
