package utils

import "testing"

func Test_CurrentFunctionName(t *testing.T) {
	expected := "github.com/puppetlabs/lumogon/utils.dummyFunction"
	actual := dummyFunction()
	if actual != expected {
		t.Errorf("Returned function name [%s] does not match expected value [%s]", actual, expected)
	}
}

func dummyFunction() string {
	return CurrentFunctionName()
}
