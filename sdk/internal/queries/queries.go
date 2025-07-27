package queries 

const LambdaTimeoutQuery = `
filter @message like /Status: timeout/
| stats count() as timeoutCount
`

const LambdaMemoryUtilizationQuery = `
parse @message "Memory Size: * MB\tMax Memory Used: * MB" as memorySize, maxMemoryUsed 
| filter ispresent(memorySize) and ispresent(maxMemoryUsed) 
| display @timestamp, memorySize, maxMemoryUsed, maxMemoryUsed / memorySize as memoryUtilizationRatio
`

const LambdaColdStartRate = `
fields @message
| filter @message like /^REPORT/
| parse @message /^REPORT.*Duration: (?<duration>[0-9\.]+) ms.*(Init Duration: (?<initDuration>[0-9\.]+) ms)?/
| stats 
    count(*) as total_invocations,
    count(initDuration) as cold_starts,
    count(initDuration) * 100.0 / count(*) as coldStartRate
`

const LambdaErrorCount = `
filter @message like /(?i)(ERROR)/
| stats count_distinct(@requestId) as errorCount
`

const LambdaUniqueRequests = `
stats count_distinct(@requestId) as invocationsCount
`

const LambdaErrorTypesQuery = `
filter @message =~ /(?i)\[ERROR\]/
| parse @message /\[ERROR\]\s+(?<error_type>[^ :]+)[ :]*(?<error_details>.*)?/
| parse error_details /(?<specific_error>.*?)\s*when calling/
| parse error_details /An error occurred \((?<aws_error_code>\w+)\)/
| stats 
    count() as error_count
    by coalesce(aws_error_code, specific_error, error_type, "UnknownError") as error_category
| sort error_count desc
`

