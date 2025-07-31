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
