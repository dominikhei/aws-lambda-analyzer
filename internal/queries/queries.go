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

package queries

const LambdaTimeoutQueryWithVersion = `
filter @type = "REPORT" and @message like /Status: timeout/ and @logStream like /\[%s\]/
| stats count_distinct(@requestId) as timeoutCount
`

const LambdaMemoryUtilizationQueryWithVersion = `
filter @type = "REPORT" and @message like /Memory Size/ and @logStream like /\[%s\]/
| parse @message "Memory Size: * MB\tMax Memory Used: * MB" as memorySize, maxMemoryUsed
| display @timestamp, memorySize, maxMemoryUsed, maxMemoryUsed / memorySize as memoryUtilizationRatio
`

const LambdaDurationQueryWithVersion = `
fields @timestamp, @message
| filter @type = "REPORT" and @message like /Duration:/ and @logStream like /\[%s\]/
| parse @message "Duration: * ms" as durationMs
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
filter @type = "REPORT" and @logStream like /\[%s\]/
| stats count_distinct(@requestId) as invocationsCount
`

const LambdaErrorTypesQueryWithVersion = `
filter @logStream like /\[%s\]/ and @message like /(?i)\[ERROR\]/
| parse @message "[ERROR] *: *" as error_type, error_details
| parse error_details "* when calling *" as specific_error, _
| parse error_details /An error occurred \((?<aws_error_code>\w+)\)/
| stats
    count() as error_count
    by coalesce(aws_error_code, specific_error, error_type, "UnknownError") as error_category
| sort error_count desc
`

const LambdaBilledDurationQueryWithVersion = `
filter @type = "REPORT" and @logStream like /\[%s\]/
| stats sum(@duration) as totalDuration, sum(@billedDuration) as totalBilledDuration
`

const LambdaColdStartDurationQueryWithVersion = `
fields @timestamp, @message
| filter @type = "REPORT" and @message like /Init Duration/ and @logStream like /\[%s\]/
| parse @message "Init Duration: * ms" as coldStartDurationMs
| filter ispresent(coldStartDurationMs)
`
