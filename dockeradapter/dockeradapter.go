package dockeradapter

import (
	"context"
	"io"
	"os"
	"strconv"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/types"
	"github.com/puppetlabs/lumogon/utils"
)

// MinSupportedAPIVersion is the lowest Docker API version that Lumogon supports
// Docker API Version 1.21 - Docker Engine 1.10.x
// - for support < 1.21 need to be able to identify the scheduler container so it can be excluded from results
// - for support < 1.20 need to use an alternative to copy when detecting the target containers OS
const MinSupportedAPIVersion = "1.21"

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
	Diff
	ServerAPIVersion() string
}

// Harvester interface exposes methods used by Capabilties Harvest functions
type Harvester interface {
	Inspector
	Executor
	Diff
	HostInspector
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
	ServerVersion(ctx context.Context) (dockertypes.Version, error)
}

// ImageInspectorPuller interface exposes methods required to both pull and
// and inspect an image
type ImageInspectorPuller interface {
	ImageInspector
	ImagePuller
}

// Creator interface exposes methods required to create an attached container
type Creator interface {
	ContainerCreate(ctx context.Context, command []string, envvars []string, image string, binds []string, links []string, kernelCapabilities []string, pidMode string, containerName string, autoRemove bool, labels map[string]string) (dockercontainer.ContainerCreateCreatedBody, error)
}

// Remover interface exposes methods required to remove a container
type Remover interface {
	ContainerRemove(ctx context.Context, containerID string, force bool) error
	CleanupHarvesters(ctx context.Context, key string) error
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

// Diff interface exposes methods required to determine files in running container
// that have been changed/added/removed relative to the containers image
type Diff interface {
	ContainerDiff(ctx context.Context, containerID string) ([]types.ChangedFile, error)
}

// containerLogOptions type contains values used to control logs returned
// from a container
type containerLogOptions struct {
	ShowStdout bool
	Tail       string
}

// concreteDockerClient wraps the upstream Docker API Client
type concreteDockerClient struct {
	Client     *client.Client
	APIVersion string
}

// ImageExists returns true if the imageName exists
func ImageExists(ctx context.Context, client ImageInspector, imageName string) bool {
	_, err := client.ImageInspect(ctx, imageName)
	if err != nil {
		return false
	}
	return true
}

// New returns a client connected with the highest API version supported by both the Lumogon client
// and the Docker runtime
func New() (Client, error) {

	host, _, _ := DockerEnvvars()

	serverAPIVersion, _, err := ServerInfo(host)
	if err != nil {
		logging.Debug("[Docker Adapter] Unable to determine Docker Server API version: %s", err.Error())
		return nil, err
	}

	// Connect using the API version that the server supports
	client, err := client.NewClient(host, serverAPIVersion, nil, nil)
	if err != nil {
		logging.Debug("[Docker Adapter] Unable to connect to Docker server")
		return nil, err
	}

	logging.Debug("[Docker Adapter] Creating container runtime client: Docker")
	concreteClient := new(concreteDockerClient)
	concreteClient.APIVersion = serverAPIVersion
	concreteClient.Client = client

	return concreteClient, nil
}

// DockerEnvvars returns values for the following environment variables, setting a default
// if no variable is set: DOCKER_HOST, DOCKER_CERT_PATH, DOCKER_TLS_VERIFY
func DockerEnvvars() (string, string, bool) {
	logging.Debug("[Docker Adapter] Negotiating client connection with Docker server")
	host, ok := os.LookupEnv("DOCKER_HOST")
	if !ok {
		host = "unix:///var/run/docker.sock"
	}
	logging.Debug("[Docker Adapter] Setting Docker host to: %s", host)

	certPath, ok := os.LookupEnv("DOCKER_CERT_PATH")
	if !ok {
		certPath = ""
	}
	logging.Debug("[Docker Adapter] Setting cert path to: %s", certPath)

	var tlsVerify bool
	s, ok := os.LookupEnv("DOCKER_TLS_VERIFY")
	if !ok {
		tlsVerify = false
	} else {
		var parseBoolErr error
		tlsVerify, parseBoolErr = strconv.ParseBool(s)
		if parseBoolErr != nil {
			logging.Debug("Error parsing DOCKER_TLS_VERIFY envvar, setting to false: %s", parseBoolErr.Error())
		}
	}
	logging.Debug("[Docker Adapter] Setting tlsVerify to: %t", tlsVerify)

	return host, certPath, tlsVerify
}

// ServerInfo returns the Server APIVersion and the servers ID
func ServerInfo(host string) (string, string, error) {
	client, err := client.NewClient(host, "", nil, nil)
	if err != nil {
		return "", "", err
	}

	ctx := context.Background()
	version, err := client.ServerVersion(ctx)
	if err != nil {
		return "", "", err
	}
	resp, _ := client.Info(ctx)

	if err = client.Close(); err != nil {
		return "", "", err
	}

	return version.APIVersion, resp.ID, nil
}

// ImageInspect inspects that requested image
func (c *concreteDockerClient) ServerAPIVersion() string {
	return c.APIVersion
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

// ServerVersion returns the underlying Docker Version struct exposed via the Engine API
func (c *concreteDockerClient) ServerVersion(ctx context.Context) (dockertypes.Version, error) {
	resp, err := c.Client.ServerVersion(ctx)
	if err != nil {
		return dockertypes.Version{}, err
	}

	return resp, nil
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
func (c *concreteDockerClient) ContainerCreate(ctx context.Context, command []string, envvars []string, image string, binds []string, links []string, kernelCapabilities []string, pidMode string, containerName string, autoRemove bool, labels map[string]string) (dockercontainer.ContainerCreateCreatedBody, error) {
	config := dockercontainer.Config{
		Image:  image,
		Cmd:    command,
		Env:    envvars,
		Labels: labels,
	}

	attachPid := dockercontainer.PidMode(pidMode)
	hostConfig := dockercontainer.HostConfig{
		CapAdd:     kernelCapabilities,
		PidMode:    attachPid,
		Binds:      binds,
		Links:      links,
		AutoRemove: autoRemove,
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
// Note that the file returned is within a tarball.
func (c *concreteDockerClient) CopyFromContainer(ctx context.Context, container, srcPath string, followSymlink bool) (io.ReadCloser, dockertypes.ContainerPathStat, error) {
	readCloser, containerPathStat, err := c.Client.CopyFromContainer(ctx, container, srcPath)
	if followSymlink && err == nil && containerPathStat.LinkTarget != "" {
		logging.Debug("[Docker Adapter] Resolving symlink for: %s, to: %s", srcPath, containerPathStat.LinkTarget)
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
		logging.Debug("[Docker Adapter] Error listing running containers: %s", err)
		return nil, err
	}

	for _, container := range containers {
		result = append(result, container.ID)
	}

	return result, nil
}

// ContainerDiff returns a slice of changed files
func (c *concreteDockerClient) ContainerDiff(ctx context.Context, containerID string) ([]types.ChangedFile, error) {
	result := []types.ChangedFile{}

	diffs, err := c.Client.ContainerDiff(ctx, containerID)
	if err != nil {
		logging.Debug("[Docker Adapter] Error getting ContainerDiff: %s", err)
		return nil, err
	}

	for _, diff := range diffs {
		result = append(result, types.ChangedFile{
			Kind: diff.Kind,
			Path: diff.Path,
		})
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

func (c *concreteDockerClient) CleanupHarvesters(ctx context.Context, key string) error {
	logging.Debug("[CleanupHarvesters] Forcefully removing attached harvester containers with the key %s", key)
	filters := filters.NewArgs()
	filters.Add("label", key)

	listOptions := dockertypes.ContainerListOptions{
		All:     true,
		Filters: filters,
	}

	containers, err := c.Client.ContainerList(ctx, listOptions)
	if err != nil {
		return err
	}

	if len(containers) == 0 {
		logging.Debug("[CleanupHarvesters] No attached harvester containers found.")
		return nil
	}

	logging.Debug("[CleanupHarvesters] Found %d attached harvester containers:", len(containers))
	for _, container := range containers {
		logging.Debug("[CleanupHarvesters] - ID: %s, State: %s", container.ID, container.State)
	}

	removeOptions := dockertypes.ContainerRemoveOptions{
		Force: true,
	}
	for _, container := range containers {
		logging.Debug("[CleanupHarvesters] Removing container ID: %s.", container.ID)
		err := c.Client.ContainerRemove(ctx, container.ID, removeOptions)
		if err != nil {
			logging.Debug("Error removing attached harvester container %s: %s\n", container.ID, err.Error())
		}
	}
	return nil
}
