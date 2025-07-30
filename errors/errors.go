package errors

import "fmt"

// The NoInvocationsError is a custom error that is thrown when a lambda function has not been invoked
// in the specified interval. This is to let users handle this special case easier, e.g set metrics to Na or 0.
type NoInvocationsError struct {
	FunctionName string
}

func (e *NoInvocationsError) Error() string {
	return fmt.Sprintf("function %q has zero invocations", e.FunctionName)
}

func NewNoInvocationsError(functionName string) error {
	return &NoInvocationsError{FunctionName: functionName}
}
