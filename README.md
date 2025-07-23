<img width="600" height="250" alt="Serverless Statistics Header" src="https://github.com/user-attachments/assets/b0711a33-a48d-45ad-a8c7-ab44e52a2286" />

Serverless Statistics is a Go sdk that allows you to extract various statistics on your AWS lambda functions, that can be used for monitoring, performance analysis, improvements and automations. It provides a simple and clean interface and handles everything for you such that you can just focus on upstream logic and the data will be provided for you. For examples on how to use it refer to the examples section.

### Available Metrics
- Cold Start Latency (avg. p95, p99)
- Memory usage (avg. p95, p99)
- num. timeout errors 
- Availability (avg. p95, p99)
- General Latency (avg. p95, p99)
- DLQ delivery success rate 
- Used / Allocated Memory
- cold start rate 
- error type breakdown (JSON)
- billed vs actual duration
- Throttles / Total Invocations
