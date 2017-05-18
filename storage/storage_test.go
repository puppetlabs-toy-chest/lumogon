package storage

import (
	"testing"

	"github.com/puppetlabs/lumogon/types"
)

var formatReportTests = []struct {
	title       string
	report      types.Report
	indent      bool
	expected    []byte
	expectError bool
}{
	{
		title:  "Empty report returns indented base JSON result",
		report: types.Report{},
		indent: true,
		expected: []byte(`{
  "$schema": "",
  "generated": "",
  "owner": "",
  "group": null,
  "client_version": {
    "BuildVersion": "",
    "BuildTime": "",
    "BuildSHA": ""
  },
  "reportid": "",
  "containers": null
}`),
		expectError: false,
	},
	{
		title:       "Empty report returns unindented base JSON result",
		report:      types.Report{},
		indent:      false,
		expected:    []byte(`{"$schema":"","generated":"","owner":"","group":null,"client_version":{"BuildVersion":"","BuildTime":"","BuildSHA":""},"reportid":"","containers":null}`),
		expectError: false,
	},
	// Note: while it is possible to test indentation on more complete report data,
	// this makes these tests more brittle than they already are -- requiring updates
	// whenever the Report type and its included types change.
}

func Test_formatReport(t *testing.T) {
	for _, test := range formatReportTests {
		t.Run(test.title, func(t *testing.T) {
			actual, err := formatReport(test.report, test.indent)
			if err != nil {
				if !test.expectError {
					t.Errorf("Unexpected error thrown: %s", err)
				}
			}

			if string(actual[:]) != string(test.expected[:]) {
				t.Errorf("Test [%s] formatReport test failed, returned result [%s], does not match expected result [%s]",
					test.title,
					actual,
					test.expected,
				)
			}
		})
	}
}
