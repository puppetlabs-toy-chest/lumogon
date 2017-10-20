package ospackages

import (
	"bufio"
	"bytes"
	"context"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/puppetlabs/lumogon/test/mocks"
)

func successfulContainerExecCreateFn(ctx context.Context, containerID string, cmd []string, attachStdout bool, attachStderr bool) (dockertypes.IDResponse, error) {
	mockIDResponse := dockertypes.IDResponse{}
	mockIDResponse.ID = "mockid"
	return mockIDResponse, nil
}

func createSuccesfulContainerExecAttachFn(buf []byte) func(context.Context, string) (dockertypes.HijackedResponse, error) {
	var succesfulContainerExecAttachFn = func(ctx context.Context, execID string) (dockertypes.HijackedResponse, error) {
		mockHijackedResponse := dockertypes.HijackedResponse{}
		mockHijackedResponse.Conn = mocks.MockNetConn{
			CloseFn: func() error {
				return nil
			},
		}
		mockReader := bytes.NewReader(buf)
		mockHijackedResponse.Reader = bufio.NewReader(mockReader)

		return mockHijackedResponse, nil
	}
	return succesfulContainerExecAttachFn
}

func successfulContainerExecStartFn(ctx context.Context, execID string) error {
	return nil
}

func successfulContainerExecInspectFn(ctx context.Context, execID string) (dockertypes.ContainerExecInspect, error) {
	mockContainerExecInspect := dockertypes.ContainerExecInspect{}
	mockContainerExecInspect.ExitCode = 0
	return mockContainerExecInspect, nil
}
