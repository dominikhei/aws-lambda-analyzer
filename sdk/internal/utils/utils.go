package utils

import (
	"context"
	"errors"
	"fmt"
	"math"
	"slices"
	"sort"
	"strings"

	"gonum.org/v1/gonum/stat"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	sdktypes "github.com/dominikhei/aws-lambda-analyzer/sdk/types"
)

type summaryStatistics struct {
	Mean      float32
	Median    float32
	P99       *float32 //Pointer as these values are nilable if the sample size is too small
	P95       *float32
	ConfInt95 *float32
	Min       float32
	Max       float32
}

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

func  AddQualifierFilter(query, qualifier string) string {
    if qualifier == "" || qualifier == "$LATEST" {
        return query
    }
    
    patterns := []string{
        fmt.Sprintf("Version: %s", qualifier),
        fmt.Sprintf("Alias: %s", qualifier),
        qualifier, 
    }
    
    var filters []string
    for _, pattern := range patterns {
        filters = append(filters, fmt.Sprintf("@message like /%s/", pattern))
    }
    
    qualifierFilter := fmt.Sprintf("filter (%s)", strings.Join(filters, " or "))
    return fmt.Sprintf("%s\n| %s", qualifierFilter, query)
}

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

func QualifierExists(ctx context.Context, client *lambda.Client, functionName, qualifier string) (bool, error) {
    if qualifier == "" || qualifier == "$LATEST" {
        return true, nil
    }
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