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

package utils

import (
	"context"
	"errors"
	"fmt"
	"math"
	"slices"
	"sort"

	"gonum.org/v1/gonum/stat"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"

	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

// summaryStatistics holds common descriptive statistics for a sample set of float64 values.
// P95, P99, and ConfInt95 are pointers because they may be nil if sample size is insufficient.
type summaryStatistics struct {
	Mean      float64
	Median    float64
	P99       *float64
	P95       *float64
	ConfInt95 *float64
	Min       float64
	Max       float64
}

// ToLoadOptions converts ConfigOptions into AWS SDK config.LoadOptions functional options.
// This abstraction should simplify configuration for users.
func ToLoadOptions(opts sdktypes.ConfigOptions) ([]func(*config.LoadOptions) error, error) {
	var loadOptions []func(*config.LoadOptions) error

	if opts.Profile != "" {
		loadOptions = append(loadOptions, func(lo *config.LoadOptions) error {
			lo.SharedConfigProfile = opts.Profile
			return nil
		})
	}

	if opts.Region != "" {
		loadOptions = append(loadOptions, func(lo *config.LoadOptions) error {
			lo.Region = opts.Region
			return nil
		})
	}

	if opts.AccessKeyID != "" && opts.SecretAccessKey != "" {
		creds := credentials.NewStaticCredentialsProvider(opts.AccessKeyID, opts.SecretAccessKey, "")
		loadOptions = append(loadOptions, func(lo *config.LoadOptions) error {
			lo.Credentials = creds
			return nil
		})
	} else if (opts.AccessKeyID != "" && opts.SecretAccessKey == "") || (opts.AccessKeyID == "" && opts.SecretAccessKey != "") {
		return nil, fmt.Errorf("both AccessKeyID and SecretAccessKey must be set together")
	}

	return loadOptions, nil
}

// CalcSummaryStats calculates descriptive statistics (mean, median, percentiles, confidence interval, min, max)
// from a slice of float64 values. Returns an error if the input slice is empty.
func CalcSummaryStats(vals []float64) (summaryStatistics, error) {
	if len(vals) == 0 {
		return summaryStatistics{}, errors.New("empty slice")
	}

	sorted := make([]float64, len(vals))
	copy(sorted, vals)
	sort.Float64s(sorted)

	mean := stat.Mean(vals, nil)
	median := stat.Quantile(0.5, stat.Empirical, sorted, nil)
	stddev := stat.StdDev(vals, nil)
	min := slices.Min(vals)
	max := slices.Max(vals)

	var p95, p99, confInt95 *float64

	if len(vals) >= 20 {
		val := stat.Quantile(0.95, stat.Empirical, sorted, nil)
		p95 = &val
	}
	if len(vals) >= 100 {
		val := stat.Quantile(0.99, stat.Empirical, sorted, nil)
		p99 = &val
	}
	if len(vals) >= 30 {
		val := 1.96 * stddev / math.Sqrt(float64(len(vals)))
		confInt95 = &val
	}

	return summaryStatistics{
		Mean:      mean,
		Median:    median,
		P95:       p95,
		P99:       p99,
		ConfInt95: confInt95,
		Min:       min,
		Max:       max,
	}, nil
}

// FunctionExists checks if an AWS Lambda function with the given name exists in the AWS account.
// Returns true if the function exists, false if not found, or an error on other failures.
func FunctionExists(ctx context.Context, client *lambda.Client, functionName string) (bool, error) {
	_, err := client.GetFunction(ctx, &lambda.GetFunctionInput{
		FunctionName: aws.String(functionName),
	})

	var nfe *types.ResourceNotFoundException
	if errors.As(err, &nfe) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// QualifierExists checks if a specific qualifier (version or alias) exists for an AWS Lambda function.
// Returns true if the qualifier exists, false if not found, or an error on other failures.
func QualifierExists(ctx context.Context, client *lambda.Client, functionName, qualifier string) (bool, error) {
	_, err := client.GetFunction(ctx, &lambda.GetFunctionInput{
		FunctionName: aws.String(functionName),
		Qualifier:    aws.String(qualifier),
	})
	var nfe *types.ResourceNotFoundException
	if errors.As(err, &nfe) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
