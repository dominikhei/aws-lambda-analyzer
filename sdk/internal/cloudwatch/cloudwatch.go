package cloudwatchfetcher

import (
    "context"
    "fmt"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/cloudwatch"
    "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
    "github.com/dominikhei/aws-lambda-analyzer/sdk/internal/clientmanager"
    sdktypes "github.com/dominikhei/aws-lambda-analyzer/sdk/types"
)

type Fetcher struct {
    client *cloudwatch.Client
}

func New(clients *clientmanager.AWSClients) *Fetcher {
    return &Fetcher{
        client: clients.CloudWatchClient,
    }
}

// FetchMetric fetches metric data from CloudWatch for a given function and metric name/stat.
func (f *Fetcher) FetchMetric(
    ctx context.Context,
    query sdktypes.FunctionQuery,
    metricName string,
    stat string,
    period int32,
) ([]types.MetricDataResult, error) {
    dimensions := []types.Dimension{
        {
            Name:  aws.String("FunctionName"),
            Value: aws.String(query.FunctionName),
        },
    }
    if query.Qualifier != "" {
        dimensions = append(dimensions, types.Dimension{
            Name:  aws.String("Resource"),
            Value: aws.String(fmt.Sprintf("%s:%s", query.FunctionName, query.Qualifier)),
        })
    }

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