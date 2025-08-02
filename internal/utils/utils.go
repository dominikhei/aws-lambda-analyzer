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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	cwTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"

	sdkinterfaces "github.com/dominikhei/serverless-statistics/internal/interfaces"
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

// mean calculates the arithmetic mean of a slice of float64 values
func mean(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}

// stdDev calculates the standard deviation of a slice of float64 values
func stdDev(vals []float64) float64 {
	if len(vals) <= 1 {
		return 0
	}

	m := mean(vals)
	sumSquares := 0.0
	for _, v := range vals {
		diff := v - m
		sumSquares += diff * diff
	}

	// Population standard deviation (n) as we consider the whole population
	// within the interval
	return math.Sqrt(sumSquares / float64(len(vals)))
}

// quantile calculates the quantile (0.0 to 1.0) from a sorted slice
func quantile(p float64, sorted []float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	index := int(math.Ceil(p*float64(len(sorted)))) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(sorted) {
		index = len(sorted) - 1
	}
	return sorted[index]
}

// CalcSummaryStats calculates descriptive statistics without external dependencies
func CalcSummaryStats(vals []float64) (summaryStatistics, error) {
	if len(vals) == 0 {
		return summaryStatistics{}, errors.New("empty slice")
	}

	sorted := make([]float64, len(vals))
	copy(sorted, vals)
	sort.Float64s(sorted)

	meanVal := mean(vals)
	medianVal := quantile(0.5, sorted)
	stddevVal := stdDev(vals)
	min := slices.Min(vals)
	max := slices.Max(vals)

	var p95, p99, confInt95 *float64

	if len(vals) >= 20 {
		val := quantile(0.95, sorted)
		p95 = &val
	}

	if len(vals) >= 100 {
		val := quantile(0.99, sorted)
		p99 = &val
	}

	if len(vals) >= 30 {
		val := 1.96 * stddevVal / math.Sqrt(float64(len(vals)))
		confInt95 = &val
	}

	return summaryStatistics{
		Mean:      meanVal,
		Median:    medianVal,
		P95:       p95,
		P99:       p99,
		ConfInt95: confInt95,
		Min:       min,
		Max:       max,
	}, nil
}

// FunctionExists checks if an AWS Lambda function with the given name exists in the AWS account.
// Returns true if the function exists, false if not found, or an error on other failures.
func FunctionExists(ctx context.Context, client sdkinterfaces.LambdaClient, functionName string) (bool, error) {
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

// QualifierExists checks if a specific qualifier (version) exists for an AWS Lambda function. Aliases are not supported.
// Returns true if the qualifier exists, false if not found, or an error on other failures.
func QualifierExists(ctx context.Context, client sdkinterfaces.LambdaClient, functionName, qualifier string) (bool, error) {
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

func SumMetricValues(results []cwTypes.MetricDataResult) (float64, error) {
	var sum float64
	for _, result := range results {
		for _, val := range result.Values {
			sum += val
		}
	}
	return sum, nil
}
