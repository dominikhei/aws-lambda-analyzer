package utils 

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	sdktypes "github.com/dominikhei/aws-lambda-analyzer/sdk/types"
)

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