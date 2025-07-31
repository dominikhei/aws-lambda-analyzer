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

package tests

import (
	"testing"

	"github.com/dominikhei/serverless-statistics/errors"
)

func TestNoInvocationsError_Error(t *testing.T) {
	err := &errors.NoInvocationsError{FunctionName: "my-test-function"}

	expected := `function "my-test-function" has zero invocations`
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}
