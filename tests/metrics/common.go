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

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"

	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

// The common file holds the mocks used by all test cases of the metrics package.

// Mock CloudWatchFetcher based on the interface in the interfaces package.
type mockCWFetcher struct {
	results []types.MetricDataResult
	err     error
}

func (m *mockCWFetcher) FetchMetric(ctx context.Context, query sdktypes.FunctionQuery, metricName string, stat string) ([]types.MetricDataResult, error) {
	return m.results, m.err
}

// Mock LogsInsights based on the interface in the interfaces package.
type mockLogsFetcher struct {
	results []map[string]string
	err     error
}

func (m *mockLogsFetcher) RunQuery(ctx context.Context, fq sdktypes.FunctionQuery, queryString string) ([]map[string]string, error) {
	return m.results, m.err
}

// Mock Lambda client based on the interface in the interfaces package.
type mockLambdaClient struct {
	GetFunctionFunc func(ctx context.Context, params *lambda.GetFunctionInput, optFns ...func(*lambda.Options)) (*lambda.GetFunctionOutput, error)
}

func (m *mockLambdaClient) GetFunction(ctx context.Context, params *lambda.GetFunctionInput, optFns ...func(*lambda.Options)) (*lambda.GetFunctionOutput, error) {
	return m.GetFunctionFunc(ctx, params, optFns...)
}
