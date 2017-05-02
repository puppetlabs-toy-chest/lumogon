package cmd

import "testing"

func TestNoGoconvey_renderVersionTemplate(t *testing.T) {
	expectedVersionOutput := `Client:
 Version:      testversionstring
 Git commit:   testbuildsha
 Built:        testdatestring`

	renderedVersionOutput, _ := renderVersionTemplate()

	if renderedVersionOutput != expectedVersionOutput {
		t.Errorf("Rendered version template [%s] does not match expected value [%s]", renderedVersionOutput, expectedVersionOutput)
	}
}
