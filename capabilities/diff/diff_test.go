package diff

import (
	"context"
	"reflect"
	"testing"

	"fmt"

	"github.com/docker/docker/pkg/archive"
	"github.com/puppetlabs/lumogon/test/mocks"
	"github.com/puppetlabs/lumogon/types"
)

// Uses generic successful mock functions from diffhelpers_test.go
var diffCapabilityTests = []struct {
	title            string
	mockDockerClient *mocks.MockDockerClient
	expected         map[string]interface{}
	expectError      bool
}{
	{
		title: "Success",
		mockDockerClient: &mocks.MockDockerClient{
			ContainerDiffFn: createContainerDiffFn(
				[]types.ChangedFile{
					types.ChangedFile{Kind: archive.ChangeAdd, Path: "/this/is/a/new/file"},
					types.ChangedFile{Kind: archive.ChangeDelete, Path: "/this/is/a/deleted/file"},
					types.ChangedFile{Kind: archive.ChangeModify, Path: "/this/is/a/modified/file"},
				},
				nil,
			),
		},
		expected: map[string]interface{}{
			"/this/is/a/new/file":      "Added",
			"/this/is/a/deleted/file":  "Deleted",
			"/this/is/a/modified/file": "Modified",
		},
		expectError: false,
	},
	{
		title: "Success No Changed Files",
		mockDockerClient: &mocks.MockDockerClient{
			ContainerDiffFn: createContainerDiffFn(
				[]types.ChangedFile{},
				nil,
			),
		},
		expected:    map[string]interface{}{},
		expectError: false,
	},
	{
		title: "Error getting container diff",
		mockDockerClient: &mocks.MockDockerClient{
			ContainerDiffFn: createContainerDiffFn(
				nil,
				fmt.Errorf("Error throw getting container diff"),
			),
		},
		expected:    map[string]interface{}{},
		expectError: true,
	},
}

func Test_diffCapability(t *testing.T) {
	for _, test := range diffCapabilityTests {
		t.Run(test.title, func(t *testing.T) {

			id := "dummyharvestid"
			target := types.TargetContainer{ID: "dummyid", Name: "dummyname", OSID: "dummyosid"}
			ctx := context.Background()

			actual, err := getChangedFiles(ctx, test.mockDockerClient, id, target)
			if err != nil {
				t.Logf("Error thrown %v", err)
				if !test.expectError {
					t.Errorf("Unexpected error thrown: %s", err)
				}
			}
			if err == nil && !reflect.DeepEqual(actual, test.expected) {
				// TODO - reflect deepEqual doesn't like empty vars
				if !(len(actual) == 0 && len(test.expected) == 0) {
					t.Errorf("Output [%v][type: %s], does not match expected results [%v][type: %s]",
						actual,
						reflect.TypeOf(actual),
						test.expected,
						reflect.TypeOf(test.expected),
					)
				}
			}
		})
	}
}
