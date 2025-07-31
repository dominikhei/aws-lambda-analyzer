package tests

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"

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
