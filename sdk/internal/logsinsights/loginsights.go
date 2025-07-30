package logsinsightsfetcher

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cloudwatchlogstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	sdktypes "github.com/dominikhei/aws-lambda-analyzer/sdk/types"
)

type Fetcher struct {
	client *cloudwatchlogs.Client
}

func New(clients *sdktypes.AWSClients) *Fetcher {
	return &Fetcher{client: clients.LogsClient}
}

func (f *Fetcher) RunQuery(ctx context.Context, fq sdktypes.FunctionQuery, queryString string) ([]map[string]string, error) {
	logGroup := fmt.Sprintf("/aws/lambda/%s", fq.FunctionName)

	startResp, err := f.client.StartQuery(ctx, &cloudwatchlogs.StartQueryInput{
		LogGroupNames: []string{logGroup},
		QueryString:   aws.String(queryString),
		StartTime:     aws.Int64(fq.StartTime.Unix()),
		EndTime:       aws.Int64(fq.EndTime.Unix()),
	})
	if err != nil {
		return nil, err
	}

	if startResp.QueryId == nil {
		return nil, errors.New("no query ID returned")
	}

	queryID := startResp.QueryId

	for {
		time.Sleep(2 * time.Second)
		resp, err := f.client.GetQueryResults(ctx, &cloudwatchlogs.GetQueryResultsInput{
			QueryId: queryID,
		})
		if err != nil {
			return nil, err
		}

		switch resp.Status {
		case cloudwatchlogstypes.QueryStatusComplete:
			var results []map[string]string
			for _, row := range resp.Results {
				m := map[string]string{}
				for _, f := range row {
					m[*f.Field] = *f.Value
				}
				results = append(results, m)
			}
			return results, nil
		case cloudwatchlogstypes.QueryStatusFailed, cloudwatchlogstypes.QueryStatusCancelled:
			return nil, fmt.Errorf("query failed with status: %s", resp.Status)
		}
	}
}
