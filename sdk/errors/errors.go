package errors 

import "fmt"

type NoInvocationsError struct {
    FunctionName string
}

func (e *NoInvocationsError) Error() string {
    return fmt.Sprintf("function %q has zero invocations", e.FunctionName)
}

func NewNoInvocationsError(functionName string) error {
    return &NoInvocationsError{FunctionName: functionName}
}