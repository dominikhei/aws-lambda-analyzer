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

package logsinsightsfetcher

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cloudwatchlogstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

// Fetcher is a wrapper around the AWS CloudWatch Logs client tailored for executing
// Logs Insights queries against Lambda function log groups.
type Fetcher struct {
	client *cloudwatchlogs.Client
}

func New(clients *sdktypes.AWSClients) *Fetcher {
	return &Fetcher{client: clients.LogsClient}
}

// RunQuery executes a Logs Insights query on the log group of the specified Lambda function,
// scoped to the time range in the FunctionQuery.
//
// Parameters:
//   - ctx: context for cancellation and timeout.
//   - fq: sdktypes.FunctionQuery containing the Lambda FunctionName, Qualifier (unused here, as it is plugged into the query before),
//     StartTime, and EndTime defining the query time window.
//   - queryString: the Logs Insights query string to execute.
//
// Returns:
//   - A slice of maps representing the query results, where each map corresponds to a row,
//     mapping field names to string values.
//   - An error if the query fails to start, returns no query ID, or fails/cancels during execution.
//
// Behavior:
//   - The function constructs the log group name using the Lambda function name in the standard
//     `/aws/lambda/{functionName}` format.
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
	// There is a 10s max duration to a query before it cancels
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	for {
		// The status is polled in a loop every 500MS.
		time.Sleep(500 * time.Millisecond)
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
		// Check if the 10s timeout is already exceeded.
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("query polling timed out")
		default:
		}
	}
}
