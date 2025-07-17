package queries 

const LambdaTimeoutQuery = `
filter @message like /Status: timeout/
| stats count() as timeoutCount
`

