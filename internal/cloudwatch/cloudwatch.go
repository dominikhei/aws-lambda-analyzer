package cloudwatchfetcher

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

// Fetcher is a wrapper around the AWS CloudWatch client tailored to fetch
// Lambda metrics efficiently using predefined dimensions and query parameters.
type Fetcher struct {
	client *cloudwatch.Client
}

// period is a default for the period parameter of Cloudwatch Metrics.
// It is set as a default and not-changeable by the user as the functionality
// does not require aggregation of metrics over sub-periods.
const period int32 = 86400

func New(clients *sdktypes.AWSClients) *Fetcher {
	return &Fetcher{
		client: clients.CloudWatchClient,
	}
}

// FetchMetric fetches metric data for a given Lambda function within the specified
// time range.
//
// Parameters:
//   - ctx: context for cancellation and deadlines.
//   - query: FunctionQuery struct containing FunctionName, Qualifier, StartTime,
//     and EndTime for the metric fetch.
//   - metricName: the name of the Lambda metric to query (e.g., "Invocations").
//   - stat: the statistic to retrieve (e.g., "Sum", "Average").
//
// Returns a slice of MetricDataResult structs containing the queried metric data,
// or an error if the request fails.
func (f *Fetcher) FetchMetric(
	ctx context.Context,
	query sdktypes.FunctionQuery,
	metricName string,
	stat string,
) ([]types.MetricDataResult, error) {
	dimensions := []types.Dimension{
		{
			Name:  aws.String("FunctionName"),
			Value: aws.String(query.FunctionName),
		},
	}

	var resourceValue string
	// In this case the version will not be part of the resource only the function name.
	if query.Qualifier == "$LATEST" {
		resourceValue = query.FunctionName
	} else {
		// For any other verison tag, the Resource dimension will be name:tag.
		resourceValue = fmt.Sprintf("%s:%s", query.FunctionName, query.Qualifier)
	}

	dimensions = append(dimensions, types.Dimension{
		Name:  aws.String("Resource"),
		Value: aws.String(resourceValue),
	})

	input := &cloudwatch.GetMetricDataInput{
		StartTime: aws.Time(query.StartTime),
		EndTime:   aws.Time(query.EndTime),
		MetricDataQueries: []types.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &types.MetricStat{
					Metric: &types.Metric{
						Namespace:  aws.String("AWS/Lambda"),
						MetricName: aws.String(metricName),
						Dimensions: dimensions,
					},
					Period: aws.Int32(period),
					Stat:   aws.String(stat),
				},
				ReturnData: aws.Bool(true),
			},
		},
	}

	resp, err := f.client.GetMetricData(ctx, input)
	if err != nil {
		return nil, err
	}

	return resp.MetricDataResults, nil
}
