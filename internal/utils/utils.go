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
	Mean      float32
	Median    float32
	P99       *float32
	P95       *float32
	ConfInt95 *float32
	Min       float32
	Max       float32
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

	var p95, p99, confInt95 *float32

	if len(vals) >= 20 {
		val := float32(stat.Quantile(0.95, stat.Empirical, sorted, nil))
		p95 = &val
	}
	if len(vals) >= 100 {
		val := float32(stat.Quantile(0.99, stat.Empirical, sorted, nil))
		p99 = &val
	}
	if len(vals) >= 30 {
		val := float32(1.96 * stddev / math.Sqrt(float64(len(vals))))
		confInt95 = &val
	}

	return summaryStatistics{
		Mean:      float32(mean),
		Median:    float32(median),
		P95:       p95,
		P99:       p99,
		ConfInt95: confInt95,
		Min:       float32(min),
		Max:       float32(max),
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
