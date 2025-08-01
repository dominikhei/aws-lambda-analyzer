<img width="480" height="200" alt="Serverless Statistics Header" src="https://github.com/user-attachments/assets/a2a24da4-cf2e-42f0-ac00-17efb6e3497c" />

Serverless Statistics is a Go sdk that allows you to extract various statistics on your AWS lambda functions, that can be used for monitoring, performance analysis, improvements and automations. It provides a simple and clean interface and handles everything for you such that you can just focus on upstream logic and the data will be provided for you. For examples on how to use it refer to the examples section.

<p align="center">
  <a href="https://goreportcard.com/report/github.com/dominikhei/serverless-statistics">
    <img src="https://goreportcard.com/badge/github.com/dominikhei/serverless-statistics" alt="Go Report Card" />
  </a>
  <a href="./LICENSE">
    <img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License" />
  </a>
  <a href="https://app.codecov.io/github/dominikhei/serverless-statistics">
    <img src="https://codecov.io/github/dominikhei/serverless-statistics/graph/badge.svg?token=MXEE07C73L" alt="Codecov" />
  </a>
  </a>
  <a href="https://github.com/dominikhei/serverless-statistics/actions/workflows/ci.yml">
    <img src="https://github.com/dominikhei/serverless-statistics/actions/workflows/ci.yml/badge.svg" alt="CI" />
  </a>
</p>

## Table of Contents

- [Key Configurations](#key-configurations)
- [Input Parameters](#input-parameters)
- [Available Metrics](#available-metrics)
- [Detailed Metric Explanations](#detailed-metric-explanations)
- [Required Permissions & CloudWatch Logging](#required-permissions--cloudwatch-logging)
- [Examples](#examples)
- [Contributing](#contributing)

## Key Configurations

### Credentials
The SDK uses the standard AWS credentials chain for authentication and will look for credentials in that order. The excemption from this is, if you specify `AccessKeyID` and `SecretAccessKey` in [ConfigOptions](./sdk/types/types.go#L42-L47) in the application code.
```go
opts := types.ConfigOptions{
  AccessKeyID:  "example-key-id",
  SecretAccessKey: "example-secret-key",
  Region: "eu-central-1",
}
stats := serverlessstatistics.New(ctx, opts)
```

### Region Configuration
When creating your client, you can specify the AWS region and other options via [ConfigOptions](./sdk/types/types.go#L42-L47). The client will look for functions in either:

- The default region from your AWS credentials/profile
- A specific region you specify in the configuration

```go
opts := types.ConfigOptions{
	Region:  "eu-central-1",
}
stats := serverlessstatistics.New(ctx, opts)
```

### Function Targeting
You specify which Lambda function to analyze by providing:

- __Function Name__: The name of your Lambda function
- __Version__ (optional): A specific version. Defaults to `$LATEST`.

If no version is provided, the SDK will analyze the logs and metrics for the `$LATEST` version.

### How are Versions and Aliases considered?
If a version tag is present, only the invocations for that specific version will be considered. If no tag is set in the `Version` parameter, the `$LATEST` version will be used by default.

Note: When using `$LATEST`, if your function was updated during the specified time frame, invocations from both the old and new versions will be included in the results (since both were `$LATEST` at different times). Set your timeframe carefully to avoid mixing versions unintentionally.

There currently is no possibility to distinguish between different aliases in the query results.

### What happens when a function has not been invoked in the specified interval?
Since the goal was to let the user decide freely what to do in this case, a [custom error](./errors/errors.go) is thrown. You can use `errors.As` in your downstream logic to asses whether this error is raised and decide yourself how you want to treat this case.
Whether a function has been invoked in a specific interval is cached internally, to reduce calls to cloudwatch metrics and thus alos charges to your AWS account.

## Input Parameters

Most metric functions in `serverless-statistics` require the following input parameters:

| Parameter     | Type              | Description |
|---------------|-------------------|-------------|
| `ctx`         | `context.Context` | Go context for timeout and cancellation. Pass `context.Background()` or derive from upstream logic. |
| `functionName`| `string`          | The name of the AWS Lambda function to analyze. This must match the name used in the AWS Console. |
| `version`   | `string` (optional) | The version of the Lambda function. Defaults to `"$LATEST"` if left empty. |
| `startTime`   | `time.Time`       | Start of the time window for analysis. Should be within the function's log retention period. |
| `endTime`     | `time.Time`       | End of the time window for analysis. Typically `time.Now()`. Must be after `startTime`. |

## Available Metrics

- [Cold Start Rate](#cold-start-rate)
- [Memory Usage Statistics](#memory-usage-statistics)
- [Timeout Rate](#timeout-rate)
- [Throttle Rate](#throttle-rate)
- [Error Rate](#error-rate)
- [Error Types](#error-types)
- [Duration Statistics](#duration-statistics)
- [Waste Ratio](#waste-ratio)
- [Cold Start Duration Statistics](#cold-start-duration-statistics)



## Detailed Metric Explanations

### Cold Start Rate

- **Source**: Logs Insights
- **Formula**:
  `count(invocations with initDuration) / count(all invocations)`
- **Return Type**: `float64`
- **Description**:
  Measures the proportion of Lambda invocations that experienced a cold start, where the Lambda execution environment had to be initialized before handling the request.

---

### Memory Usage Statistics

- **Source**: Logs Insights
- **Formula**:
  `Maximum Used Memory / Allocated Memory`
- **Return Type**: `float64`
- **Available Aggregations**:
  - Minimum Memory Usage Rate
  - Maximum Memory Usage Rate
  - Median Memory Usage Rate
  - Mean Memory Usage Rate
  - 95th Percentile Memory Usage Rate (requires ≥ 20 invocations)
  - 99th Percentile Memory Usage Rate (requires ≥ 100 invocations)
  - 95% Confidence Interval of Memory Usage Rate (requires ≥ 30 invocations)
- **Description**:
  Indicates how much of the allocated memory the Lambda function actually uses during execution.
- **Notes**:
  For each invocation, the peak memory is considered.
---

### Throttle Rate

- **Source**: CloudWatch
- **Formula**:
    `count(invocations with throttling) / count(all invocations)`
- **Return Type**: `float64`
- **Description**:
  Shows the fraction of Lambda invocations that were throttled due to exceeding concurrency limits.

---

### Timeout Rate

- **Source**: CloudWatch
- **Formula**:
  `count(invocations with timeout) / count(all invocations)`
- **Return Type**: `float64`
- **Description**:
  Measures how often Lambda functions exceed their configured execution timeout.

---

### Error Rate

- **Source**: Logs Insights
- **Formula**:
  `count(invocations with error) / count(all invocations)`
- **Return Type**: `float64`
- **Description**:
  Represents the proportion of Lambda invocations that resulted in errors, identified by presence of error logs.
- **Notes**:
  Based on distinct `requestID`s where `[ERROR]` is present in logs.

---

### Error Types

- **Source**: Logs Insights
- **Return Type**: `[]ErrorType`
- **Description**:
  Categorizes errors seen in logs by their semantic type (e.g., `ImportError`, `TimeoutError`).
- **Notes**:
  Function timeouts do not count as errors
---
### Duration Statistics

- **Source**: Logs Insights
- **Return Type**: `float64`
- **Available Aggregations**:
  - Minimum Duration
  - Maximum Duration
  - Median Duration
  - Mean Duration
  - 95th Percentile Duration (requires ≥ 20 invocations)
  - 99th Percentile Duration (requires ≥ 100 invocations)
  - 95% Confidence Interval of Duration (requires ≥ 30 invocations)
- **Description**:
  Provides detailed timing metrics of Lambda execution duration.
---

### Waste Ratio

- **Source**: Logs Insights
- **Formula**:
  `(billed duration - execution duration) / billed duration`
- **Return Type**: `float64`
- **Description**:
  Measures the fraction of billed time during which the function was not actively executing code (e.g., waiting or initialization overhead).
---

### Cold Start Duration Statistics

- **Source**: Logs Insights
- **Return Type**: `float64`
- **Available Aggregations**:
  - Minimum Duration
  - Maximum Duration
  - Median Duration
  - Mean Duration
  - 95th Percentile Duration (requires ≥ 20 invocations)
  - 99th Percentile Duration (requires ≥ 100 invocations)
  - 95% Confidence Interval of Duration (requires ≥ 30 invocations)
- **Description**:
  Provides statistics on the time spent initializing initializing the Lambda execution environment during cold starts.
---

### Function Configuration

- **Source**: Lambda API
- **Available Aggregations**:
  - Function Name
  - Function ARN
  - Function Version
  - Configured Memory Size
  - Configured Timeout
  - Runtime
  - Last Modification Date
  - Environment Variables
- **Description**:
  Retrieves the current configuration settings of a Lambda function.
---

## Required Permissions & CloudWatch Logging

To successfully retrieve and analyze Lambda metrics and logs, the SDK requires your AWS credentials to have specific IAM permissions, and the target Lambda functions must have CloudWatch logging enabled.

### Required IAM Permissions

The following minimum IAM permissions are required to use the SDK:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:StartQuery",
        "logs:GetQueryResults",
        "logs:DescribeLogGroups",
        "logs:DescribeLogStreams",
        "logs:FilterLogEvents",
        "cloudwatch:GetMetricData",
        "cloudwatch:GetMetricStatistics",
        "lambda:GetFunctionConfiguration",
        "lambda:ListFunctions"
      ],
      "Resource": "*"
    }
  ]
}
```

## Examples

This section aims to provide two short example on how to use the sdk. For the available methods and metrics refer to [Available Metrics](#available-metrics) and the detailed section on each metric below. You simply need to pass in the function name, version, starttime and endtime into a metric function. The return will be a custom struct exposed in [types](./types/types.go). There you can also see how to access the relevant values and which ones are available.

### Initialization:

```go
import (
	serverlessstatistics "github.com/dominikhei/serverless-statistics"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

ctx := context.Background()
opts := sdktypes.ConfigOptions{
	Region:  "eu-central-1",
	Profile: "default",
}
stats := serverlessstatistics.New(ctx, opts)
```

### Memory Usage Statistics:

```go
stats := serverlessstatistics.New(ctx, opts)
functionName := "testFunction"
version := "v1"
endTime := time.Now()
startTime := endTime.Add(-12 * time.Hour)
rate, err := stats.GetWasteRatio(ctx, functionName, version, startTime, endTime)
if err != nil {
	fmt.Printf("error: %v\n", err)
	return
}
fmt.Printf("The waste ratio is %v percent", rate.WasteRatio * 100)

```

## Contributing

If you find any bugs or have ideas for new functionality, please open an Issue or, even better, submit a Pull Request. All contributions are welcome!
