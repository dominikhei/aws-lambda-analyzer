<img width="1200" height="500" alt="Serverless Statistics Header" src="https://github.com/user-attachments/assets/2fd91fe3-74c0-4b39-b60a-f56dfa0e57c4" />

Serverless Statistics is a Go sdk that allows you to extract various statistics on your AWS lambda functions, that can be used for monitoring, performance analysis, improvements and automations. It provides a simple and clean interface and handles everything for you such that you can just focus on upstream logic and the data will be provided for you. For examples on how to use it refer to the examples section.

## How are Versions and Aliases considered?
If a version tag is present, only the invocations for that specific version will be considered. If no tag is set in the `qualifier` parameter, the `$LATEST` qualifier will be used by default.

Note: When using `$LATEST`, if your function was updated during the specified time frame, invocations from both the old and new versions will be included in the results (since both were `$LATEST` at different times). Set your timeframe carefully to avoid mixing versions unintentionally.

There currently is no possibility to distinguish between different aliases in the query results.

## Key Configurations

### Credentials
The SDK uses the standard AWS credentials chain for authentication and will look for credentials in that order. The excemption from this is, if you specify `AccessKeyID` and `SecretAccessKey` in [ConfigOptions](./sdk/types/types.go#L42-L47) in the application code.

### Region Configuration
When creating your client, you can specify the AWS region and other options via [ConfigOptions](./sdk/types/types.go#L42-L47). The client will operate in either:

- The default region from your AWS credentials/profile
- A specific region you specify in the configuration

### Function Targeting
You specify which Lambda function to analyze by providing:

- __Function Name__ - The name of your Lambda function
- __Qualifier__ (optional) - A specific version number. Defaults to `$LATEST`.

If no qualifier is provided, the SDK will analyze all logs and metrics for the function across all versions.


## Available Metrics

- [Cold Start Rate](#cold-start-rate)
- [Memory Usage Statistics](#memory-usage-statistics)
- [Timeout Rate](#timeout-rate)
- [Throttle Rate](#throttle-rate)
- [Error Rate](#error-rate)
- [Error Types](#error-types)



## Detailed Metric Explanations

### Cold Start Rate

- **Source**: Logs Insights  
- **Formula**:  
  `sum(invocations with initDuration) / sum(all invocations)`  
- **Return Type**: `float32`

---

### Memory Usage Statistics

- **Source**: Logs Insights  
- **Formula**:  
  `Used Memory / Allocated Memory`  
- **Return Type**: `float32`  
- **Available Aggregations**:  
  - Minimum Memory Usage Rate  
  - Maximum Memory Usage Rate  
  - Median Memory Usage Rate  
  - Mean Memory Usage Rate  
  - 95th Percentile Memory Usage Rate (requires ≥ 20 invocations)  
  - 99th Percentile Memory Usage Rate (requires ≥ 100 invocations)  
  - 95% Confidence Interval of Memory Usage Rate (requires ≥ 30 invocations)

---

### Throttle Rate

- **Source**: CloudWatch  
- **Formula**:  
  `Throttles / Invocations`  
- **Return Type**: `float32`

---

### Timeout Rate

- **Source**: CloudWatch  
- **Formula**:  
  `Timeouts / Invocations`  
- **Return Type**: `float32`

---

### Error Rate

- **Source**: Logs Insights  
- **Formula**:  
  `sum(invocations with error) / sum(all invocations)`  
- **Return Type**: `float32`  
- **Notes**:  
  Based on distinct `requestID`s where `[ERROR]` is present in logs.

---

### Error Types

- **Source**: Logs Insights  
- **Return Type**: `[]ErrorType`  
- **Description**:  
  Extracted from logs containing `[ERROR]`, grouped by semantic category (e.g., timeout, dependency error).  
- **Notes**:  
  Function timeouts do not count as errors