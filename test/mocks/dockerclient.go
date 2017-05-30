package mocks

import (
	"context"
	"fmt"
	"io"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/puppetlabs/lumogon/types"
	"github.com/puppetlabs/lumogon/utils"
)

// MockDockerClient is a mock type implementing the dockeradapter.Client interface, by default
// all interface functions will raise a panic, you can however provide your own test specific
// function implementations which will take precedence. Refer to the dockerruntime_test.go
// file for examples
type MockDockerClient struct {
	Type                   string
	ImagePullFn            func(ctx context.Context, imageName string) error
	ContainerInspectFn     func(ctx context.Context, containerID string) (dockertypes.ContainerJSON, error)
	ContainerCreateFn      func(ctx context.Context, command []string, envvars []string, image string, binds []string, links []string, kernelCapabilities []string, pidMode string, containerName string) (dockercontainer.ContainerCreateCreatedBody, error)
	ImageExistsFn          func(ctx context.Context, imageName string) bool
	ContainerRemoveFn      func(ctx context.Context, containerID string, force bool) error
	ContainerStartFn       func(ctx context.Context, containerID string) error
	ContainerLogsFn        func(ctx context.Context, containerID string) (string, error)
	ContainerListFn        func(ctx context.Context) ([]string, error)
	ContainerExecCreateFn  func(ctx context.Context, containerID string, cmd []string, attachStdout bool, attachStderr bool) (dockertypes.IDResponse, error)
	ContainerExecStartFn   func(ctx context.Context, execID string) error
	ContainerExecAttachFn  func(ctx context.Context, execID string, cmd []string, attachStdout bool, attachStderr bool) (dockertypes.HijackedResponse, error)
	ContainerExecInspectFn func(ctx context.Context, execID string) (dockertypes.ContainerExecInspect, error)
	ImageInspectFn         func(ctx context.Context, imageName string) (dockertypes.ImageInspect, error)
	CopyFromContainerFn    func(ctx context.Context, container, srcPath string, followSymlink bool) (io.ReadCloser, dockertypes.ContainerPathStat, error)
	ContainerFilesystemFn  func(ctx context.Context, containerID string) (types.Filesystem, error)
}

// ImagePull is a mock implementation of dockeradapter.ImagePull
func (c MockDockerClient) ImagePull(ctx context.Context, imageName string) error {
	if c.ImagePullFn != nil {
		fmt.Println("[MockDockerClient] In ", utils.CurrentFunctionName())
		fmt.Println("[MockDockerClient]  - ctx: ", ctx)
		fmt.Println("[MockDockerClient]  - imageName: ", imageName)
		return c.ImagePullFn(ctx, imageName)
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// ContainerInspect is a mock implementation of dockeradapter.ContainerInspect
func (c MockDockerClient) ContainerInspect(ctx context.Context, containerID string) (dockertypes.ContainerJSON, error) {
	if c.ContainerInspectFn != nil {
		fmt.Println("[MockDockerClient] In ", utils.CurrentFunctionName())
		fmt.Println("[MockDockerClient]  - ctx: ", ctx)
		fmt.Println("[MockDockerClient]  - containerID: ", containerID)
		return c.ContainerInspectFn(ctx, containerID)
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// ContainerCreate is a mock implementation of dockeradapter.ContainerCreate
func (c MockDockerClient) ContainerCreate(ctx context.Context, command []string, envvars []string, image string, binds []string, links []string, kernelCapabilities []string, pidMode string, containerName string) (dockercontainer.ContainerCreateCreatedBody, error) {
	if c.ContainerCreateFn != nil {
		fmt.Println("[MockDockerClient] In ", utils.CurrentFunctionName())
		fmt.Println("[MockDockerClient]  - ctx: ", ctx)
		fmt.Println("[MockDockerClient]  - command: ", command)
		fmt.Println("[MockDockerClient]  - envvars: ", envvars)
		fmt.Println("[MockDockerClient]  - image: ", envvars)
		fmt.Println("[MockDockerClient]  - binds: ", binds)
		fmt.Println("[MockDockerClient]  - links: ", links)
		fmt.Println("[MockDockerClient]  - kernelCapabilities: ", kernelCapabilities)
		fmt.Println("[MockDockerClient]  - pidMode: ", pidMode)
		fmt.Println("[MockDockerClient]  - containerName: ", containerName)
		return c.ContainerCreateFn(ctx, command, envvars, image, binds, links, kernelCapabilities, pidMode, containerName)
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// ImageExists is a mock implementation of dockeradapter.ImageExists
func (c MockDockerClient) ImageExists(ctx context.Context, imageName string) bool {
	if c.ImageExistsFn != nil {
		fmt.Println("[MockDockerClient] In ", utils.CurrentFunctionName())
		fmt.Println("[MockDockerClient]  - ctx: ", ctx)
		fmt.Println("[MockDockerClient]  - imageName: ", imageName)
		return c.ImageExistsFn(ctx, imageName)
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// ContainerRemove is a mock implementation of dockeradapter.ContainerRemove
func (c MockDockerClient) ContainerRemove(ctx context.Context, containerID string, force bool) error {
	if c.ContainerRemoveFn != nil {
		fmt.Println("[MockDockerClient] In ", utils.CurrentFunctionName())
		fmt.Println("[MockDockerClient]  - ctx: ", ctx)
		fmt.Println("[MockDockerClient]  - containerID: ", containerID)
		fmt.Println("[MockDockerClient]  - force: ", force)
		return c.ContainerRemoveFn(ctx, containerID, force)
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// ContainerStart is a mock implementation of dockeradapter.ContainerStart
func (c MockDockerClient) ContainerStart(ctx context.Context, containerID string) error {
	if c.ContainerStartFn != nil {
		fmt.Println("[MockDockerClient] In ", utils.CurrentFunctionName())
		fmt.Println("[MockDockerClient]  - ctx: ", ctx)
		fmt.Println("[MockDockerClient]  - containerID: ", containerID)
		return c.ContainerStartFn(ctx, containerID)
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// ContainerLogs is a mock implementation of dockeradapter.ContainerLogs
func (c MockDockerClient) ContainerLogs(ctx context.Context, containerID string) (string, error) {
	if c.ContainerLogsFn != nil {
		fmt.Println("[MockDockerClient] In ", utils.CurrentFunctionName())
		fmt.Println("[MockDockerClient]  - ctx: ", ctx)
		fmt.Println("[MockDockerClient]  - containerID: ", containerID)
		return c.ContainerLogsFn(ctx, containerID)
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// ContainerList is a mock implementation of dockeradapter.ContainerList
func (c MockDockerClient) ContainerList(ctx context.Context) ([]string, error) {
	if c.ContainerLogsFn != nil {
		fmt.Println("[MockDockerClient] In ", utils.CurrentFunctionName())
		fmt.Println("[MockDockerClient]  - ctx: ", ctx)
		return c.ContainerListFn(ctx)
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// ContainerExecCreate is a mock implementation of dockeradapter.ContainerExecCreate
func (c *MockDockerClient) ContainerExecCreate(ctx context.Context, containerID string, cmd []string, attachStdout bool, attachStderr bool) (dockertypes.IDResponse, error) {
	if c.ContainerExecCreateFn != nil {
		fmt.Println("[MockDockerClient] In ", utils.CurrentFunctionName())
		fmt.Println("[MockDockerClient]  - ctx: ", ctx)
		fmt.Println("[MockDockerClient]  - containerID: ", containerID)
		fmt.Println("[MockDockerClient]  - cmd: ", cmd)
		fmt.Println("[MockDockerClient]  - attachStdout: ", attachStdout)
		fmt.Println("[MockDockerClient]  - attachStderr: ", attachStderr)
		return c.ContainerExecCreateFn(ctx, containerID, cmd, attachStdout, attachStderr)
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// ContainerExecStart is a mock implementation of dockeradapter.ContainerExecStart
func (c *MockDockerClient) ContainerExecStart(ctx context.Context, execID string) error {
	if c.ContainerExecStartFn != nil {
		fmt.Println("[MockDockerClient] In ", utils.CurrentFunctionName())
		fmt.Println("[MockDockerClient]  - ctx: ", ctx)
		fmt.Println("[MockDockerClient]  - execID: ", execID)
		return c.ContainerExecStartFn(ctx, execID)
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// ContainerExecAttach is a mock implementation of dockeradapter.ContainerExecAttach
func (c *MockDockerClient) ContainerExecAttach(ctx context.Context, execID string, cmd []string, attachStdout bool, attachStderr bool) (dockertypes.HijackedResponse, error) {
	if c.ContainerExecAttachFn != nil {
		fmt.Println("[MockDockerClient] In ", utils.CurrentFunctionName())
		fmt.Println("[MockDockerClient]  - ctx: ", ctx)
		fmt.Println("[MockDockerClient]  - execID: ", execID)
		fmt.Println("[MockDockerClient]  - cmd: ", cmd)
		fmt.Println("[MockDockerClient]  - attachStdout: ", attachStdout)
		fmt.Println("[MockDockerClient]  - attachStderr: ", attachStderr)
		return c.ContainerExecAttachFn(ctx, execID, cmd, attachStdout, attachStderr)
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// ContainerExecInspect is a mock implementation of dockeradapter.ContainerExecInspect
func (c *MockDockerClient) ContainerExecInspect(ctx context.Context, execID string) (dockertypes.ContainerExecInspect, error) {
	if c.ContainerExecInspectFn != nil {
		fmt.Println("[MockDockerClient] In ", utils.CurrentFunctionName())
		fmt.Println("[MockDockerClient]  - ctx: ", ctx)
		fmt.Println("[MockDockerClient]  - execID: ", execID)
		return c.ContainerExecInspectFn(ctx, execID)
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// ImageInspect is a mock implementation of dockeradapter.ImageInspect
func (c MockDockerClient) ImageInspect(ctx context.Context, imageName string) (dockertypes.ImageInspect, error) {
	if c.ImageInspectFn != nil {
		fmt.Println("[MockDockerClient] In ", utils.CurrentFunctionName())
		fmt.Println("[MockDockerClient]  - ctx: ", ctx)
		fmt.Println("[MockDockerClient]  - imageName: ", imageName)
		return c.ImageInspectFn(ctx, imageName)
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// CopyFromContainer is a mock implementation of dockeradapter.CopyFromContainer
func (c MockDockerClient) CopyFromContainer(ctx context.Context, container, srcPath string, followSymlink bool) (io.ReadCloser, dockertypes.ContainerPathStat, error) {
	if c.CopyFromContainerFn != nil {
		fmt.Println("[MockDockerClient] In ", utils.CurrentFunctionName())
		return c.CopyFromContainerFn(ctx, container, srcPath, followSymlink)
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// ContainerFilesystem is a mock implementation of dockeradapter.ContainerFilesystem
func (c MockDockerClient) ContainerFilesystem(ctx context.Context, containerID string) (types.Filesystem, error) {
	if c.ContainerFilesystemFn != nil {
		fmt.Println("[MockDockerClient] In ", utils.CurrentFunctionName())
		return c.ContainerFilesystemFn(ctx, containerID)
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}
