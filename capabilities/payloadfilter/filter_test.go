package payloadfilter

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type filterTest struct {
	description string
	input       map[string]interface{}
	expected    map[string]interface{}
	expectError bool
}

var filterTests = []filterTest{
	{
		"An empty map should result in an empty map",
		make(map[string]interface{}),
		make(map[string]interface{}),
		false,
	},
	{
		"A valid map should be returned unchanged",
		map[string]interface{}{
			"a": "1",
			"b": "2",
			"c": "3",
			"d": "4",
		},
		map[string]interface{}{
			"a": "1",
			"b": "2",
			"c": "3",
			"d": "4",
		},
		false,
	},
	{
		"Native JSON ints should be passed through unchanged",
		map[string]interface{}{
			"a": 1,
			"b": 2,
			"c": "3",
			"d": "4",
		},
		map[string]interface{}{
			"a": 1,
			"b": 2,
			"c": "3",
			"d": "4"},
		false,
	},
	{
		"A map with complex non-string values should fail with an error",
		map[string]interface{}{
			"a": "1",
			"b": map[string][]string{
				"a": {"1", "2", "3"},
				"b": {"1", "2", "3"},
			},
			"c": "3",
			"d": "4",
		},
		map[string]interface{}{},
		true,
	},
	{
		"A valid map with empty values should drop those pairs",
		map[string]interface{}{
			"a": "",
			"b": "2",
			"c": "3",
			"d": "4",
		},
		map[string]interface{}{
			"b": "2",
			"c": "3",
			"d": "4",
		},
		false,
	},
	{
		"A valid map with empty keys should drop those pairs",
		map[string]interface{}{
			"":  "1",
			"b": "2",
			"c": "3",
			"d": "4",
		},
		map[string]interface{}{
			"b": "2",
			"c": "3",
			"d": "4",
		},
		false,
	},
	{
		"A valid map with a nested map of strings should be returned unchanged",
		map[string]interface{}{
			"strings": map[string]string{"a": "1", "b": "2"},
			"b":       "2",
			"c":       "3",
			"d":       "4",
		},
		map[string]interface{}{
			"strings": map[string]string{"a": "1", "b": "2"},
			"b":       "2",
			"c":       "3",
			"d":       "4",
		},
		false,
	},
	{
		"A valid map with an empty map of strings should drop the empty map of strings",
		map[string]interface{}{
			"strings": map[string]string{},
			"b":       "2",
			"c":       "3",
			"d":       "4",
		},
		map[string]interface{}{
			"b": "2",
			"c": "3",
			"d": "4",
		},
		false,
	},
	{
		"A valid map with nested maps with empty string values should drop empty value pairs",
		map[string]interface{}{
			"strings": map[string]string{"a": "1", "b": "", "c": "3"},
			"b":       "2",
			"c":       "3",
			"d":       "4",
		},
		map[string]interface{}{
			"strings": map[string]string{"a": "1", "c": "3"},
			"b":       "2",
			"c":       "3",
			"d":       "4",
		},
		false,
	},
	{
		"A valid map with nested maps with empty key values should drop empty value pairs",
		map[string]interface{}{
			"strings": map[string]string{"": "1", "b": "2", "c": "3", "d": ""},
			"b":       "2",
			"c":       "3",
			"d":       "4",
		},
		map[string]interface{}{
			"strings": map[string]string{"b": "2", "c": "3"},
			"b":       "2",
			"c":       "3",
			"d":       "4",
		},
		false,
	},
	{
		"A deep structure with lots of embedded empty values should be cleaned up properly",
		map[string]interface{}{
			"strings": map[string]string{"a": "1", "": "2", "c": "3", "d": ""},
			"empty":   "",
			"maps on maps on maps": map[string]interface{}{
				"z": "26",
				"y": "x",
				"w": map[string]string{"a": "1", "": "2", "c": "3", "d": ""},
				"v": map[string]string{},
				"u": map[string]interface{}{
					"t": "",
					"s": map[string]string{
						"r": "",
					},
					"q": "p",
				},
			},
			"b": "2",
			"c": "3",
			"d": "4",
		},
		map[string]interface{}{
			"strings": map[string]string{"a": "1", "c": "3"},
			"maps on maps on maps": map[string]interface{}{
				"z": "26",
				"y": "x",
				"w": map[string]string{"a": "1", "c": "3"},
				"u": map[string]interface{}{
					"q": "p",
				},
			},
			"b": "2",
			"c": "3",
			"d": "4",
		},
		false,
	},
}

func TestFilter(t *testing.T) {
	for _, example := range filterTests {
		actual, err := Filter(example.input)
		if err != nil {
			if !example.expectError {
				t.Errorf("Error [%s] when running Filter for example '%s' on input:\n\t%s\n",
					err,
					example.description,
					spew.Sprintf("%#+v", example.input))
			}
			continue
		}

		if !reflect.DeepEqual(actual, example.expected) {
			t.Errorf("\nFilter test failed: %s\nexpected:\n\t%s\nactual:\n\t%s\n",
				example.description,
				spew.Sprintf("%#+v", example.expected),
				spew.Sprintf("%#+v", actual))
		}
	}
}
