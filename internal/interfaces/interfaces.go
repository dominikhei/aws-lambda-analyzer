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

package fetcherinterfaces

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/dominikhei/serverless-statistics/internal/cache"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

// This interface matches logsinsightsfetcher.Fetcher for tetsing the internal functions
type LogsInsightsFetcher interface {
	RunQuery(ctx context.Context, fq sdktypes.FunctionQuery, queryString string) ([]map[string]string, error)
}

// This interface matches cloudwatchfetcher.Fetcher for tetsing the internal functions
type CloudWatchFetcher interface {
	FetchMetric(ctx context.Context, query sdktypes.FunctionQuery, metricName string, stat string) ([]types.MetricDataResult, error)
}

// This interface matches lambda.Client for tetsing the internal functions
type LambdaClient interface {
	GetFunction(ctx context.Context, params *lambda.GetFunctionInput, optFns ...func(*lambda.Options)) (*lambda.GetFunctionOutput, error)
}

type Cache interface {
	Has(key cache.CacheKey) bool
	Set(key cache.CacheKey, value int)
	Get(key cache.CacheKey) (int, bool)
}
