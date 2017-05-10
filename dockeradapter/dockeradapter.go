package dockeradapter

import (
	"context"
	"fmt"
	"io"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/utils"
)

// Client is a Docker (currently local) ContainerRuntime
type Client interface {
	ImagePuller
	ImageInspector
	Inspector
	Executor
	Creator
	Remover
	Starter
	LogGetter
	Lister
	HostInspector
	CopyFrom
}

// Harvester interface exposes methods used by Capabilties Harvest functions
type Harvester interface {
	Inspector
	Executor
}

// ImagePuller interface exposes methods required to pull an image
type ImagePuller interface {
	ImagePull(ctx context.Context, imageName string) error
}

// ImageInspector interface exposes methods required to inspect an image
type ImageInspector interface {
	ImageInspect(ctx context.Context, imageName string) (dockertypes.ImageInspect, error)
}

// Inspector interface exposes methods required to inspect a container
type Inspector interface {
	ContainerInspect(ctx context.Context, containerID string) (dockertypes.ContainerJSON, error)
}

// Executor TODO
type Executor interface {
	ContainerExecCreate(ctx context.Context, execID string, cmd []string, attachStdout bool, attachStderr bool) (dockertypes.IDResponse, error)
	ContainerExecStart(ctx context.Context, execID string) error
	ContainerExecAttach(ctx context.Context, execID string, cmd []string, attachStdout bool, attachStderr bool) (dockertypes.HijackedResponse, error)
	ContainerExecInspect(ctx context.Context, execID string) (dockertypes.ContainerExecInspect, error)
}

// HostInspector interface exposes methods required to inspect a docker host
type HostInspector interface {
	HostID(ctx context.Context) string
}

// ImageInspectorPuller interface exposes methods required to both pull and
// and inspect an image
type ImageInspectorPuller interface {
	ImageInspector
	ImagePuller
}

// Creator interface exposes methods required to create a container
type Creator interface {
	ContainerCreate(ctx context.Context, command []string, envvars []string, image string, binds []string, links []string, kernelCapabilities []string, pidMode string, containerName string) (dockercontainer.ContainerCreateCreatedBody, error)
}

// Remover interface exposes methods required to remove a container
type Remover interface {
	ContainerRemove(ctx context.Context, containerID string, force bool) error
}

// Starter interface exposes methods required to start a container
type Starter interface {
	ContainerStart(ctx context.Context, containerID string) error
}

// LogGetter interface exposes methods required to get logs from a container
type LogGetter interface {
	ContainerLogs(ctx context.Context, containerID string) (string, error)
}

// Lister interface exposes methods required to list containers
type Lister interface {
	ContainerList(ctx context.Context) ([]string, error)
}

// CopyFrom interface exposes methods required to copy file data from a container
type CopyFrom interface {
	CopyFromContainer(ctx context.Context, container, srcPath string, followSymlink bool) (io.ReadCloser, dockertypes.ContainerPathStat, error)
}

// containerLogOptions type contains values used to control logs returned
// from a container
type containerLogOptions struct {
	ShowStdout bool
	Tail       string
}

// concreteDockerClient wraps the upstream Docker API Client
type concreteDockerClient struct {
	Client *client.Client
}

// ImageExists returns true if the imageName exists
func ImageExists(ctx context.Context, client ImageInspector, imageName string) bool {
	_, err := client.ImageInspect(ctx, imageName)
	if err != nil {
		return false
	}
	return true
}

// New returns a client satisfying the Client interface
func New() (Client, error) {
	concreteClient := new(concreteDockerClient)
	logging.Stderr("[Docker Adapter] Creating container runtime client: Docker")
	dockerAPIClient, err := client.NewEnvClient()
	if err != nil {
		return nil, fmt.Errorf("[Docker Adapter] Unable to initialise container runtime type: Docker, error: %s", err)
	}
	concreteClient.Client = dockerAPIClient
	return concreteClient, nil
}

// ImageInspect inspects that requested image
func (c *concreteDockerClient) ImageInspect(ctx context.Context, imageName string) (dockertypes.ImageInspect, error) {
	imageInspect, _, err := c.Client.ImageInspectWithRaw(ctx, imageName)
	return imageInspect, err
}

// ImagePull pulls the requested image
func (c *concreteDockerClient) ImagePull(ctx context.Context, imageName string) error {
	_, err := c.Client.ImagePull(ctx, imageName, dockertypes.ImagePullOptions{})
	return err
}

// ContainerInspect inspects the requested container
func (c *concreteDockerClient) ContainerInspect(ctx context.Context, containerID string) (dockertypes.ContainerJSON, error) {
	return c.Client.ContainerInspect(ctx, containerID)
}

// HostID returns the Unique ID of a host generated from SSH Host Keys
func (c *concreteDockerClient) HostID(ctx context.Context) string {
	resp, _ := c.Client.Info(ctx)
	return resp.ID
}

func (c *concreteDockerClient) ContainerExecCreate(ctx context.Context, containerID string, cmd []string, attachStdout bool, attachStderr bool) (dockertypes.IDResponse, error) {
	execOpts := dockertypes.ExecConfig{
		Cmd:          cmd,
		AttachStdout: attachStdout,
		AttachStderr: attachStderr,
	}
	return c.Client.ContainerExecCreate(ctx, containerID, execOpts)
}

func (c *concreteDockerClient) ContainerExecStart(ctx context.Context, execID string) error {
	execStartOpts := dockertypes.ExecStartCheck{}
	return c.Client.ContainerExecStart(ctx, execID, execStartOpts)
}

func (c *concreteDockerClient) ContainerExecAttach(ctx context.Context, execID string, cmd []string, attachStdout bool, attachStderr bool) (dockertypes.HijackedResponse, error) {
	execOpts := dockertypes.ExecConfig{
		Cmd:          cmd,
		AttachStdout: attachStdout,
		AttachStderr: attachStderr,
	}
	return c.Client.ContainerExecAttach(ctx, execID, execOpts)
}

func (c *concreteDockerClient) ContainerExecInspect(ctx context.Context, execID string) (dockertypes.ContainerExecInspect, error) {
	return c.Client.ContainerExecInspect(ctx, execID)
}

// ContainerCreate creates a container with the supplied subset of the Docker
// API configuration
func (c *concreteDockerClient) ContainerCreate(ctx context.Context, command []string, envvars []string, image string, binds []string, links []string, kernelCapabilities []string, pidMode string, containerName string) (dockercontainer.ContainerCreateCreatedBody, error) {
	config := dockercontainer.Config{
		Image: image,
		Cmd:   command,
		Env:   envvars,
	}

	attachPid := dockercontainer.PidMode(pidMode)
	hostConfig := dockercontainer.HostConfig{
		CapAdd:  kernelCapabilities,
		PidMode: attachPid,
		Binds:   binds,
		Links:   links,
	}

	return c.Client.ContainerCreate(ctx, &config, &hostConfig, nil, containerName)
}

// ContainerRemove removes the requested container
func (c *concreteDockerClient) ContainerRemove(ctx context.Context, containerID string, force bool) error {
	containerRemoveOptions := dockertypes.ContainerRemoveOptions{
		Force: force,
	}

	return c.Client.ContainerRemove(ctx, containerID, containerRemoveOptions)
}

// ContainerStart starts the requested container
func (c *concreteDockerClient) ContainerStart(ctx context.Context, containerID string) error {
	containerStartOptions := dockertypes.ContainerStartOptions{}
	return c.Client.ContainerStart(ctx, containerID, containerStartOptions)
}

// ContainerLogs (DEPRECATED) returns a string container the last line of output logged to Stdout
// from the requested container (this was used to capture the harvest output)
func (c *concreteDockerClient) ContainerLogs(ctx context.Context, containerID string) (string, error) {
	containerLogsOptions := dockertypes.ContainerLogsOptions{
		ShowStdout: true,
		Tail:       "1",
	}

	readCloser, err := c.Client.ContainerLogs(ctx, containerID, containerLogsOptions)
	if err != nil {
		return "", err
	}

	// TODO - currently stripping the header from the log output, this should go away when we move to an alternate
	// mechanism for communicating between harvester and scheduler
	logs, err := utils.GetStringFromReader(readCloser)
	if err != nil {
		return "", err
	}
	return stripDockerLogsHeader(logs), nil
}

// CopyFromContainer returns a ReadCloser containing the copied file and a
// ContainerPathStat with the files attributes. Optionally follow symlinks.
// Note that the file returned is withing a tarball.
func (c *concreteDockerClient) CopyFromContainer(ctx context.Context, container, srcPath string, followSymlink bool) (io.ReadCloser, dockertypes.ContainerPathStat, error) {
	readCloser, containerPathStat, err := c.Client.CopyFromContainer(ctx, container, srcPath)
	if followSymlink && err == nil && containerPathStat.LinkTarget != "" {
		logging.Stderr("[Docker Adapter] Resolving symlink for: %s, to: %s", srcPath, containerPathStat.LinkTarget)
		readCloser, containerPathStat, err = c.Client.CopyFromContainer(ctx, container, containerPathStat.LinkTarget)
	}
	return readCloser, containerPathStat, err
}

// ContainerList returns a slice of container ID strings
func (c *concreteDockerClient) ContainerList(ctx context.Context) ([]string, error) {
	result := []string{}
	containerListOptions := dockertypes.ContainerListOptions{}

	containers, err := c.Client.ContainerList(ctx, containerListOptions)
	if err != nil {
		logging.Stderr("[Docker Adapter] Error listing running containers: %s", err)
		return nil, err
	}

	for _, container := range containers {
		result = append(result, container.ID)
	}

	return result, nil
}

// stripDockerLogsHeader (DEPRECATED) strips the Docker logs header and aims to return a string
func stripDockerLogsHeader(rawlogs string) string {
	headerLength := 8
	if len(rawlogs) <= 8 {
		return ""
	}
	return rawlogs[headerLength:]
}
