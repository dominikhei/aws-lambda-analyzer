package queries 

const LambdaTimeoutQueryWithVersion = `
filter @message like /Status: timeout/ and @logStream like /\[%s\]/
| stats count_distinct(@requestId) as timeoutCount
`


const LambdaMemoryUtilizationQueryWithVersion = `
parse @message "Memory Size: * MB\tMax Memory Used: * MB" as memorySize, maxMemoryUsed 
| filter ispresent(memorySize) and ispresent(maxMemoryUsed)  and @logStream like /\[%s\]/
| display @timestamp, memorySize, maxMemoryUsed, maxMemoryUsed / memorySize as memoryUtilizationRatio
`

const LambdaColdStartRateWithVersion = `
filter @type = "REPORT" and @logStream like /\[%s\]/
| parse @message /REPORT RequestId: (?<requestId>[a-f0-9-]+)/
| stats 
    count_distinct(requestId) as totalInvocations,
    sum(strcontains(@message, "Init Duration")) as coldStartLines
`

const LambdaErrorCountWithVersion = `
filter @message like /(?i)(ERROR)/ and @logStream like /\[%s\]/
| stats count_distinct(@requestId) as errorCount
`

const LambdaUniqueRequestsWithVersion = `
filter @logStream like /\[%s\]/
| stats count_distinct(@requestId) as invocationsCount
`

const LambdaErrorTypesQueryWithVersion = `
filter @message =~ /(?i)\[ERROR\]/ and @logStream like /\[%s\]/
| parse @message /\[ERROR\]\s+(?<error_type>[^ :]+)[ :]*(?<error_details>.*)?/
| parse error_details /(?<specific_error>.*?)\s*when calling/
| parse error_details /An error occurred \((?<aws_error_code>\w+)\)/
| stats 
    count() as error_count
    by coalesce(aws_error_code, specific_error, error_type, "UnknownError") as error_category
| sort error_count desc
`

