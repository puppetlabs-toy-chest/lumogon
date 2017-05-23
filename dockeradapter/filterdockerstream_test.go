package dockeradapter

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/puppetlabs/lumogon/test/helper"
)

var filterDockerStreamTests = []struct {
	title            string
	input            [][]byte
	filterStreamType int
	expected         []string
	expectError      bool
}{
	{
		title: "Single line with matching stream type in payload",
		input: [][]byte{
			helper.AddDockerStreamHeader([]byte("abc\n"), 1),
		},
		filterStreamType: 1,
		expected: []string{
			"abc",
		},
		expectError: false,
	},
	{
		title: "Multiple lines, 1 with matching stream type in payload",
		input: [][]byte{
			helper.AddDockerStreamHeader([]byte("abc\n"), 2),
			helper.AddDockerStreamHeader([]byte("def\n"), 1),
		},
		filterStreamType: 1,
		expected: []string{
			"def",
		},
		expectError: false,
	},
	{
		title: "Multiple lines, all with matching stream type in payload",
		input: [][]byte{
			helper.AddDockerStreamHeader([]byte("abc\n"), 2),
			helper.AddDockerStreamHeader([]byte("def\n"), 2),
		},
		filterStreamType: 2,
		expected: []string{
			"abc",
			"def",
		},
		expectError: false,
	},
	{
		title: "Multiple lines, none with matching stream type in payload",
		input: [][]byte{
			helper.AddDockerStreamHeader([]byte("abc\n"), 0),
			helper.AddDockerStreamHeader([]byte("def\n"), 1),
			helper.AddDockerStreamHeader([]byte("ghi\n"), 2),
		},
		filterStreamType: 3,
		expected:         []string{},
		expectError:      false,
	},
	{
		title:            "Empty buffer",
		input:            [][]byte{},
		filterStreamType: 3,
		expected:         []string{},
		expectError:      false,
	},
	{
		title: "Multi-line payload with matching stream type",
		input: [][]byte{
			helper.AddDockerStreamHeader([]byte("abc\ndef\nghi\n"), 1),
		},
		filterStreamType: 1,
		expected: []string{
			"abc",
			"def",
			"ghi",
		},
		expectError: false,
	},
	{
		title: "Combination of single line and multi-line payloads with matching stream type",
		input: [][]byte{
			helper.AddDockerStreamHeader([]byte("abc\ndef\nghi\n"), 1),
			helper.AddDockerStreamHeader([]byte("jkl\n"), 1),
		},
		filterStreamType: 1,
		expected: []string{
			"abc",
			"def",
			"ghi",
			"jkl",
		},
		expectError: false,
	},
	{
		title: "Empty lines at beginning discarded",
		input: [][]byte{
			helper.AddDockerStreamHeader([]byte("\n\n\n\n\n\nabc\ndef\n"), 1),
			helper.AddDockerStreamHeader([]byte("ghi\n"), 1),
		},
		filterStreamType: 1,
		expected: []string{
			"abc",
			"def",
			"ghi",
		},
		expectError: false,
	},
	{
		title: "Empty lines mid stream discarded",
		input: [][]byte{
			helper.AddDockerStreamHeader([]byte("abc\n\nghi\n"), 1),
			helper.AddDockerStreamHeader([]byte("jkl\n\n\n"), 1),
			helper.AddDockerStreamHeader([]byte("nmo\n"), 1),
		},
		filterStreamType: 1,
		expected: []string{
			"abc",
			"ghi",
			"jkl",
			"nmo",
		},
		expectError: false,
	},
	{
		title: "Empty lines at end discarded",
		input: [][]byte{
			helper.AddDockerStreamHeader([]byte("abc\ndef\n"), 1),
			helper.AddDockerStreamHeader([]byte("ghi\n\n\n\n\n\n\n"), 1),
		},
		filterStreamType: 1,
		expected: []string{
			"abc",
			"def",
			"ghi",
		},
		expectError: false,
	},
	{
		title: "Buffer with only blank lines",
		input: [][]byte{
			helper.AddDockerStreamHeader([]byte("\n"), 1),
			helper.AddDockerStreamHeader([]byte("\n"), 1),
			helper.AddDockerStreamHeader([]byte("\n"), 1),
			helper.AddDockerStreamHeader([]byte("\n\n\n\n\n\n\n"), 1),
		},
		filterStreamType: 1,
		expected:         []string{},
		expectError:      false,
	},
	{
		title: "Buffer with no terminating newline",
		input: [][]byte{
			helper.AddDockerStreamHeader([]byte("abc\ndef"), 1),
		},
		filterStreamType: 1,
		expected: []string{
			"abc",
			"def",
		},
		expectError: false,
	},
	{
		title: "Header with invalid size - too large, reads past end of buffer",
		input: [][]byte{
			helper.AddCustomDockerStreamHeader([]byte("abc"), 1, 10),
		},
		filterStreamType: 1,
		expected:         []string{},
		expectError:      true,
	},
	{
		// This causes the subsequent iteration to interpret the byte it lands on
		// as a stream header so will behave in an unexpected fashion depending on the
		// contents of the stream.
		// TODO - unsure at this point how defensive this needs to be implmented?
		title: "Header with invalid size, too small, doesn't capture entire payload",
		input: [][]byte{
			helper.AddCustomDockerStreamHeader([]byte("abc"), 1, -1),
			helper.AddDockerStreamHeader([]byte("def"), 1),
		},
		filterStreamType: 1,
		expected:         []string{},
		expectError:      true,
	},
}

func Test_FilterDockerStream(t *testing.T) {
	for _, test := range filterDockerStreamTests {
		t.Run(test.title, func(t *testing.T) {
			var buf []byte
			for _, entry := range test.input {
				buf = append(buf, entry...)
			}
			r := bytes.NewReader(buf)
			actual, err := FilterDockerStream(r, test.filterStreamType)
			if err != nil {
				if !test.expectError {
					t.Errorf("Unexpected error thrown: %s", err)
					t.Logf("Input buffer: %v", buf)
					t.Logf("Output string slice: %v", actual)
				}
			}
			if err != nil && !reflect.DeepEqual(actual, test.expected) {
				// TODO - reflect deepEqual doesn't like empty vars
				if !(len(actual) == 0 && len(test.expected) == 0) {
					t.Errorf("Output [%v][type: %s], does not match expected results [%v][type: %s]",
						actual,
						reflect.TypeOf(actual),
						test.expected,
						reflect.TypeOf(test.expected),
					)
					t.Logf("Input buffer: %v", buf)
					t.Logf("Output string slice: %v", actual)
				}
			}
		})
	}
}

// Each frame results in two reads, one for the header and one
// for the payload, with a final read to pick up the io.EOF
// The test buffer below requires 5 successful reads to complete.
var multiFrameStdOutBuf = [][]byte{
	helper.AddDockerStreamHeader([]byte("abc\n"), 1),
	helper.AddDockerStreamHeader([]byte("def\n"), 1),
}

var buffReadErrTests = []struct {
	title            string
	input            [][]byte
	filterStreamType int
	validReads       int
	expectedError    error
}{
	{
		title:            "Error reading first frame header",
		input:            multiFrameStdOutBuf,
		filterStreamType: 1,
		validReads:       0,
		expectedError:    helper.ErrTimeout,
	},
	{
		title:            "Error reading first frame payload",
		input:            multiFrameStdOutBuf,
		filterStreamType: 1,
		validReads:       1,
		expectedError:    helper.ErrTimeout,
	},
	{
		title:            "Error reading second frame header",
		input:            multiFrameStdOutBuf,
		filterStreamType: 1,
		validReads:       2,
		expectedError:    helper.ErrTimeout,
	},
	{
		title:            "Error reading second frame payload",
		input:            multiFrameStdOutBuf,
		filterStreamType: 1,
		validReads:       3,
		expectedError:    helper.ErrTimeout,
	},
	{
		title:            "Error reading buffer EOF",
		input:            multiFrameStdOutBuf,
		filterStreamType: 1,
		validReads:       4,
		expectedError:    helper.ErrTimeout,
	},
	{
		title:            "Successfully read buffer and EOF",
		input:            multiFrameStdOutBuf,
		filterStreamType: 1,
		validReads:       5,
		expectedError:    nil,
	},
}

func Test_BuffReadErrTests(t *testing.T) {
	for _, test := range buffReadErrTests {
		t.Run(test.title, func(t *testing.T) {
			var buf []byte
			for _, entry := range test.input {
				buf = append(buf, entry...)
			}
			r := bytes.NewReader(buf)
			tr := helper.TimeoutReader(r, test.validReads)
			_, err := FilterDockerStream(tr, test.filterStreamType)
			t.Logf("Error returned: %s", err)
			if err == nil && test.expectedError == nil {
				return
			}
			if err.Error() != helper.ErrTimeout.Error() {
				t.Errorf("Expected error not thrown: %s", err)
			}
		})
	}
}
