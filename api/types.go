package api

import "time"

type MetricsSummary struct {
    AverageDurationMs          float64
    DurationP50Ms              float64
    DurationP90Ms              float64
    DurationP99Ms              float64
    MaxMemoryMB                int64
    AverageMaxMemoryUsedMB     float64
    MemoryUsagePercent         float64
    InvocationCount            int64
    ErrorCount                 int64
    ThrottleCount              int64
    ColdStartCount             int64
    AverageColdStartDurationMs float64
    MaxConcurrentExecutions    int64
    RetryCount                 int64
    DLQSendCount               int64
}

type InvocationPattern struct {
    InvocationsPerMinute map[time.Time]int64
}

type PhaseTimings struct {
    InitDurationMs        float64
    HandlerDurationMs     float64
    CleanupDurationMs     float64
    ExternalCallDurations map[string]float64
}

type ErrorMetrics struct {
    ErrorCount    int64
    ThrottleCount int64
}

type FunctionConfig struct {
    MemorySize       int32
    TimeoutSeconds   int32
    Runtime          string
    ProvisionedConcurrency int64
}