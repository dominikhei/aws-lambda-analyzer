<img width="600" height="250" alt="Serverless Statistics Header" src="https://github.com/user-attachments/assets/b0711a33-a48d-45ad-a8c7-ab44e52a2286" />

Serverless Statistics is a Go sdk that allows you to extract various statistics on your AWS lambda functions, that can be used for monitoring, performance analysis, improvements and automations. It provides a simple and clean interface and handles everything for you such that you can just focus on upstream logic and the data will be provided for you. For examples on how to use it refer to the examples section.

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
- __Qualifier__ (optional) - A specific version number or alias (e.g., "prod", "v1")

If no qualifier is provided, the SDK will analyze all logs and metrics for the function across all versions and aliases.


## Available Metrics

- [Cold Start Rate](#cold-start-rate)
- [Memory Usage Statistics](#memory-usage-statistics)
- [Timeout Rate](#timeout-rate)
- [Throttle Rate](#throttle-rate)
- [Dead Letter Queue Statistics](#dead-letter-queue-statistics)



## Detailed explanations

### Cold Start Rate

| Source        | Calculation                                                   | Aggregation Metric | Return Type | Notes |
|---------------|---------------------------------------------------------------|---------------------|-------------|-------|
| Logs Insights | sum(invocations with initDuration) / sum(all invocations)     | Cold Start Rate     | float32     |       |


### Memory Usage Statistics

| Source        | Calculation                         | Aggregation Metric                              | Return Type | Notes                                      |
|---------------|--------------------------------------|--------------------------------------------------|-------------|--------------------------------------------|
| Logs Insights | Used Memory / Allocated Memory       | Minimum Memory Usage Rate                        | float32     |                                            |
|               |                                      | Maximum Memory Usage Rate                        | float32     |                                            |
|               |                                      | Median Memory Usage Rate                         | float32     |                                            |
|               |                                      | Mean Memory Usage Rate                           | float32     |                                            |
|               |                                      | 95th Percentile Memory Usage Rate                | float32     | Requires ≥ 20 invocations                  |
|               |                                      | 99th Percentile Memory Usage Rate                | float32     | Requires ≥ 100 invocations                 |
|               |                                      | 95% Confidence Interval of Memory Usage Rate     | float32     | Requires ≥ 30 invocations                  |

### Throttle Rate

| Source        | Calculation                                                         | Aggregation Metric | Return Type | Notes |
|---------------|----------------------------------------------------------------------|---------------------|-------------|-------|
| CloudWatch    | Throttles / Invocations                                              | Throttle Rate       | float32     |       |

### Timeout Rate

| Source        | Calculation                                                         | Aggregation Metric | Return Type | Notes |
|---------------|----------------------------------------------------------------------|---------------------|-------------|-------|
| CloudWatch    | Timeouts / Invocations                                               | Timeout Rate        | float32     |       |
