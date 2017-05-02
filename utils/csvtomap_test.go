package utils

import (
	"reflect"
	"testing"
)

var csvToMapTests = []struct {
	title       string
	input       []string
	expected    map[string]string
	expectError bool
}{
	{
		title: "Extract key value pair",
		input: []string{
			"key,val",
		},
		expected: map[string]string{
			"key": "val",
		},
		expectError: false,
	},
	{
		title: "Extract multiple key value pairs",
		input: []string{
			"key1,val1",
			"key2,val2",
			"key3,val3",
		},
		expected: map[string]string{
			"key1": "val1",
			"key2": "val2",
			"key3": "val3",
		},
		expectError: false,
	},
	{
		title: "Invalid csv with >2 columns",
		input: []string{
			"key1,val1,extra1",
		},
		expectError: true,
	},
	{
		title: "Invalid csv with 1 column",
		input: []string{
			"key1",
		},
		expectError: true,
	},
	{
		title: "Invalid csv empty key",
		input: []string{
			",val1",
		},
		expectError: true,
	},
	{
		title: "Invalid empty entry",
		input: []string{
			"",
		},
		expectError: true,
	},
	{
		title: "Invalid csv with mix of valid and invalid entries",
		input: []string{
			"key1,val1",
			"",
		},
		expectError: true,
	},
	{
		title:       "Empty string slice",
		input:       []string{},
		expected:    map[string]string{},
		expectError: false,
	},
}

func Test_CsvToMap(t *testing.T) {
	for _, test := range csvToMapTests {
		t.Run(test.title, func(t *testing.T) {
			actual, err := CsvToMap(test.input)
			if err != nil {
				if !test.expectError {
					t.Errorf("Test [%s] threw unexpected error [%s]", test.title, err)
					t.Logf("Input: %v", test.input)
				}
				return
			}
			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("Test [%s] FilterDockerStream test failed, output [%v][type: %s], does not match expected results [%v][type: %s]",
					test.title,
					actual,
					reflect.TypeOf(actual),
					test.expected,
					reflect.TypeOf(test.expected),
				)
				t.Logf("Input: %v", test.input)
				t.Logf("Actual: %v", actual)
			}
		})
	}
}
