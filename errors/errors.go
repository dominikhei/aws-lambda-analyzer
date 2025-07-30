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
