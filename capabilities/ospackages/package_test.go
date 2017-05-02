package ospackages

import (
	"context"
	"reflect"
	"testing"

	"fmt"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/puppetlabs/lumogon/test/helper"
	"github.com/puppetlabs/lumogon/test/mocks"
)

// Uses generic successful mock functions from packagehelpers_test.go

var runPackageCmdTests = []struct {
	title                     string
	mockNetConn               mocks.MockNetConn
	mockHarvesterDockerClient *mocks.MockDockerClient
	expected                  map[string]string
	expectError               bool
}{
	{
		title: "Success",
		mockHarvesterDockerClient: &mocks.MockDockerClient{
			ContainerExecCreateFn: successfulContainerExecCreateFn,
			ContainerExecAttachFn: createSuccesfulContainerExecAttachFn(
				helper.AddDockerStreamHeader([]byte("packagename,0.0.1-rc-x86_64\n"), 1)),
			ContainerExecStartFn:   successfulContainerExecStartFn,
			ContainerExecInspectFn: successfulContainerExecInspectFn,
		},
		expected: map[string]string{
			"packagename": "0.0.1-rc-x86_64",
		},
		expectError: false,
	},
	{
		title: "Error creating exec",
		mockHarvesterDockerClient: &mocks.MockDockerClient{
			ContainerExecCreateFn: func(ctx context.Context, containerID string, cmd []string, attachStdout bool, attachStderr bool) (dockertypes.IDResponse, error) {
				return dockertypes.IDResponse{}, fmt.Errorf("Dummy ContainerExecCreate error")
			},
		},
		expectError: true,
	},
	{
		title: "Error attaching to exec",
		mockHarvesterDockerClient: &mocks.MockDockerClient{
			ContainerExecCreateFn: successfulContainerExecCreateFn,
			ContainerExecAttachFn: func(ctx context.Context, execID string, cmd []string, attachStdout bool, attachStderr bool) (dockertypes.HijackedResponse, error) {
				return dockertypes.HijackedResponse{}, fmt.Errorf("Dummy ContainerExecAttach error")
			},
		},
		expectError: true,
	},
	{
		title: "Error starting exec",
		mockHarvesterDockerClient: &mocks.MockDockerClient{
			ContainerExecCreateFn: successfulContainerExecCreateFn,
			ContainerExecAttachFn: createSuccesfulContainerExecAttachFn(
				helper.AddDockerStreamHeader([]byte("packagename,0.0.1-rc-x86_64\n"), 1)),
			ContainerExecStartFn: func(ctx context.Context, execID string) error {
				return fmt.Errorf("Dummy ContainerExecStart error")
			},
		},
		expectError: true,
	},
	{
		title: "Error running ContainerExecInspect",
		mockHarvesterDockerClient: &mocks.MockDockerClient{
			ContainerExecCreateFn: successfulContainerExecCreateFn,
			ContainerExecAttachFn: createSuccesfulContainerExecAttachFn(
				helper.AddDockerStreamHeader([]byte("packagename,0.0.1-rc-x86_64\n"), 1)),
			ContainerExecStartFn: successfulContainerExecStartFn,
			ContainerExecInspectFn: func(ctx context.Context, execID string) (dockertypes.ContainerExecInspect, error) {
				return dockertypes.ContainerExecInspect{}, fmt.Errorf("Dummy ContainerExecInspect error")
			},
		},
		expectError: true,
	},
	{
		title: "Non-zero exit code running command in exec",
		mockHarvesterDockerClient: &mocks.MockDockerClient{
			ContainerExecCreateFn: successfulContainerExecCreateFn,
			ContainerExecAttachFn: createSuccesfulContainerExecAttachFn(
				helper.AddDockerStreamHeader([]byte("packagename,0.0.1-rc-x86_64\n"), 1)),
			ContainerExecStartFn: successfulContainerExecStartFn,
			ContainerExecInspectFn: func(ctx context.Context, execID string) (dockertypes.ContainerExecInspect, error) {
				mockContainerExecInspect := dockertypes.ContainerExecInspect{}
				mockContainerExecInspect.ExitCode = 127
				return mockContainerExecInspect, nil
			},
		},
		expectError: true,
	},
	{
		title: "Error filtering response from exec",
		mockHarvesterDockerClient: &mocks.MockDockerClient{
			ContainerExecCreateFn: successfulContainerExecCreateFn,
			ContainerExecAttachFn: createSuccesfulContainerExecAttachFn(
				helper.AddCustomDockerStreamHeader([]byte("packagename,0.0.1-rc-x86_64\n"), 1, 10)),
			ContainerExecStartFn:   successfulContainerExecStartFn,
			ContainerExecInspectFn: successfulContainerExecInspectFn,
		},
		expectError: true,
	},
	{
		title: "Error converting response to a map",
		mockHarvesterDockerClient: &mocks.MockDockerClient{
			ContainerExecCreateFn: successfulContainerExecCreateFn,
			ContainerExecAttachFn: createSuccesfulContainerExecAttachFn(
				helper.AddDockerStreamHeader([]byte("packagename,0.0.1-rc-x86_64,unexpected extra text\n"), 1)),
			ContainerExecStartFn:   successfulContainerExecStartFn,
			ContainerExecInspectFn: successfulContainerExecInspectFn,
		},
		expectError: true,
	},
}

func Test_runPackageCmd(t *testing.T) {
	for _, test := range runPackageCmdTests {
		t.Run(test.title, func(t *testing.T) {

			testContainerID := "testcontainerid1"
			testCmd := []string{"dummy", "command"}

			actual, err := runPackageCmd(test.mockHarvesterDockerClient, testContainerID, testCmd)
			if err != nil {
				if !test.expectError {
					t.Errorf("Unexpected error thrown: %s", err)
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
				}
			}
		})
	}
}
