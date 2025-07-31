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

// GetColdStartRate calculates the cold start rate of an AWS Lambda function
// over a specified time range and qualifier (version).
// The cold start rate is computed as cold starts divided by total invocations.
func GetColdStartRate(
	ctx context.Context,
	logsFetcher sdkinterfaces.LogsInsightsFetcher, //Interface to inject mock for unit tests
	cwFetcher sdkinterfaces.CloudWatchFetcher,
	query sdktypes.FunctionQuery,
) (*sdktypes.ColdStartRateReturn, error) {

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
	queryString := fmt.Sprintf(queries.LambdaColdStartRateWithVersion, escapedQualifier)
	results, err := logsFetcher.RunQuery(ctx, query, queryString)
	if err != nil {
		return nil, fmt.Errorf("run logs insights query: %w", err)
	}
	var coldStartRate float64
	totalStr := results[0]["totalInvocations"]
	coldStr := results[0]["coldStartLines"]

	if totalStr != "" && coldStr != "" {
		total, err1 := strconv.ParseFloat(totalStr, 64)
		cold, err2 := strconv.ParseFloat(coldStr, 64)
		if err1 != nil || err2 != nil || total == 0 {
			return nil, fmt.Errorf("invalid data from logs: total=%v, cold=%v", totalStr, coldStr)
		}
		coldStartRate = cold / total
	}
	return &sdktypes.ColdStartRateReturn{
		ColdStartRate: coldStartRate,
		FunctionName:  query.FunctionName,
		Qualifier:     query.Qualifier,
		StartTime:     query.StartTime,
		EndTime:       query.EndTime,
	}, nil
}
