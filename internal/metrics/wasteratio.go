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
	sdkinterfaces "github.com/dominikhei/serverless-statistics/internal/interfaces"
	"github.com/dominikhei/serverless-statistics/internal/queries"
	"github.com/dominikhei/serverless-statistics/internal/utils"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

// GetWasteRatio calculates the ratio between billed duration and execution duration
// of an AWS Lambda function over a specified time range and qualifier (version).
// The waste ratio is computed as: (billed duration - total duration / billed duration)
func GetWasteRatio(
	ctx context.Context,
	cwFetcher sdkinterfaces.CloudWatchFetcher,
	logsFetcher sdkinterfaces.LogsInsightsFetcher,
	query sdktypes.FunctionQuery,
) (*sdktypes.WasteRatioReturn, error) {

	invocationsResults, err := cwFetcher.FetchMetric(ctx, query, "Invocations", "Sum")
	if err != nil {
		return nil, fmt.Errorf("fetch invocations metric: %w", err)
	}
	invocationsSum, err := utils.SumMetricValues(invocationsResults)
	if err != nil {
		return nil, fmt.Errorf("parse invocations metric data: %w", err)
	}
	if invocationsSum == 0 {
		return nil, &sdkerrors.NoInvocationsError{FunctionName: query.FunctionName}
	}

	escapedQualifier := strings.ReplaceAll(query.Qualifier, "$", "\\$")
	queryString := fmt.Sprintf(queries.LambdaBilledDurationQueryWithVersion, escapedQualifier)
	results, err := logsFetcher.RunQuery(ctx, query, queryString)
	if err != nil {
		return nil, fmt.Errorf("fetch errors from logs insights: %w", err)
	}
	var totalDurationMs, totalBilledDurationMs float64
	if len(results) > 0 {
		val := results[0]["totalDuration"]
		if val != "" {
			totalDurationMs, err = strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, fmt.Errorf("parse totalDurationMs from logs: %w", err)
			}
		}
		val = results[0]["totalBilledDuration"]
		if val != "" {
			totalBilledDurationMs, err = strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, fmt.Errorf("parse totalBilledDurationMs from logs: %w", err)
			}
		}
	}

	if totalDurationMs == 0 {
		return nil, fmt.Errorf("total duration is zero, cannot calculate waste ratio")
	}

	wasteRatio := (totalBilledDurationMs - totalDurationMs) / totalBilledDurationMs

	return &sdktypes.WasteRatioReturn{
		WasteRatio:   wasteRatio,
		FunctionName: query.FunctionName,
		Qualifier:    query.Qualifier,
		StartTime:    query.StartTime,
		EndTime:      query.EndTime,
	}, nil
}
