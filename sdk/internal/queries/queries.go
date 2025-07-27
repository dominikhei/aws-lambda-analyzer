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

const LambdaErrorRate = `
filter @message like /(?i)error/ 
| stats count(*) as errorInvocations, count(*) + count(@message not like /(?i)error/) as totalInvocations
`