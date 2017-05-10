package utils

import "testing"

var keyInMapTests = []struct {
	title        string
	searchString string
	inputMap     map[string]int
	expected     bool
}{
	{
		title:        "string present as key in map",
		searchString: "debian",
		inputMap: map[string]int{
			"debian": 1,
		},
		expected: true,
	},
	{
		title:        "string not present as key in map",
		searchString: "debian",
		inputMap: map[string]int{
			"fedora": 1,
		},
		expected: false,
	},
}

func Test_KeyInMap(t *testing.T) {
	for _, test := range keyInMapTests {
		t.Run(test.title, func(t *testing.T) {
			actual := KeyInMap(test.searchString, test.inputMap)
			if actual != test.expected {
				t.Errorf("Test [%s] KeyInMap test failed, returned result [%v], does not match expected result [%v]",
					test.title,
					actual,
					test.expected,
				)
				t.Logf("SearchString: %s", test.searchString)
				t.Logf("InputMap: %v", test.inputMap)
			}
		})
	}
}
