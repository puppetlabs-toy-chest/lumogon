package dockeradapter

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"testing"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/puppetlabs/lumogon/test/mocks"
	"github.com/puppetlabs/lumogon/types"
)

// stringToTargetContainer Tests
func Test_stringToTargetContainer(t *testing.T) {
	var mockClient = mocks.MockDockerClient{
		ContainerInspectFn: func(ctx context.Context, containerID string) (dockertypes.ContainerJSON, error) {
			containerJSON := dockertypes.ContainerJSON{}
			containerJSON.ContainerJSONBase = &dockertypes.ContainerJSONBase{}
			containerJSON.ContainerJSONBase.ID = fmt.Sprintf("testID_%s", containerID)
			containerJSON.ContainerJSONBase.Name = fmt.Sprintf("testName_%s", containerID)
			containerJSON.State = &dockertypes.ContainerState{}
			containerJSON.State.Running = true
			return containerJSON, nil
		},
	}
	ctx := context.TODO()

	expectedTargetContainer := types.TargetContainer{
		ID:   "testID_1",
		Name: "testName_1",
	}
	actualTargetContainer, err := stringToTargetContainer(ctx, "1", mockClient)
	if err != nil {
		t.Errorf("Unexpected error thrown by stringToTargetContainer: %s", err)
	}
	if actualTargetContainer.ID != expectedTargetContainer.ID {
		t.Errorf("TargetContainer.ID [%s] does not match expected value [%s]", actualTargetContainer.ID, expectedTargetContainer.ID)
	}
	if actualTargetContainer.Name != expectedTargetContainer.Name {
		t.Errorf("TargetContainer.Name [%s] does not match expected value [%s]", actualTargetContainer.Name, expectedTargetContainer.Name)
	}
}
func Test_stringToTargetContainer_targetNotRunning(t *testing.T) {
	var mockClient = mocks.MockDockerClient{
		ContainerInspectFn: func(ctx context.Context, containerID string) (dockertypes.ContainerJSON, error) {
			containerJSON := dockertypes.ContainerJSON{}
			containerJSON.ContainerJSONBase = &dockertypes.ContainerJSONBase{}
			containerJSON.ContainerJSONBase.ID = fmt.Sprintf("testID_%s", containerID)
			containerJSON.ContainerJSONBase.Name = fmt.Sprintf("testName_%s", containerID)
			containerJSON.State = &dockertypes.ContainerState{}
			containerJSON.State.Running = false
			return containerJSON, nil
		},
	}
	ctx := context.TODO()

	actualTargetContainer, err := stringToTargetContainer(ctx, "1", mockClient)
	if err != nil {
		t.Errorf("Unexpected error thrown by stringToTargetContainer: %s", err)
	}
	if actualTargetContainer != nil {
		t.Errorf("Expected nil container")
	}
}

func Test_stringsToTargetContainers(t *testing.T) {
	var mockClient = mocks.MockDockerClient{
		ContainerInspectFn: func(ctx context.Context, containerID string) (dockertypes.ContainerJSON, error) {
			containerJSON := dockertypes.ContainerJSON{}
			containerJSON.ContainerJSONBase = &dockertypes.ContainerJSONBase{}
			containerJSON.ContainerJSONBase.ID = fmt.Sprintf("testID_%s", containerID)
			containerJSON.ContainerJSONBase.Name = fmt.Sprintf("testName_%s", containerID)
			containerJSON.State = &dockertypes.ContainerState{}
			containerJSON.State.Running = true
			return containerJSON, nil
		},
	}
	ctx := context.TODO()

	expectedTargetContainer0 := types.TargetContainer{
		ID:   "testID_0",
		Name: "testName_0",
	}
	expectedTargetContainer1 := types.TargetContainer{
		ID:   "testID_1",
		Name: "testName_1",
	}
	targetContainerIDs := []string{"0", "1"}
	actualTargetContainers := stringsToTargetContainers(ctx, targetContainerIDs, mockClient)

	if actualTargetContainers[0].ID != expectedTargetContainer0.ID {
		t.Errorf("TargetContainer[0].ID [%s] does not match expected value [%s]", actualTargetContainers[0].ID, expectedTargetContainer0.ID)
	}
	if actualTargetContainers[1].Name != expectedTargetContainer1.Name {
		t.Errorf("TargetContainer[1].Name [%s] does not match expected value [%s]", actualTargetContainers[1].Name, expectedTargetContainer1.Name)
	}
}

func Test_stringsToTargetContainers_NoValidIDs(t *testing.T) {
	var mockClient = mocks.MockDockerClient{
		ContainerInspectFn: func(ctx context.Context, containerID string) (dockertypes.ContainerJSON, error) {
			return dockertypes.ContainerJSON{}, fmt.Errorf("DummyError")
		},
	}
	ctx := context.TODO()

	targetContainerIDs := []string{"0", "1"}
	actualTargetContainers := stringsToTargetContainers(ctx, targetContainerIDs, mockClient)

	if len(actualTargetContainers) != 0 {
		t.Errorf("Expected 0 TargetContainers, received %d.", len(actualTargetContainers))
	}
}

// getContainerOS Tests
var getContainerOSTests = []struct {
	title                 string
	containerID           string
	mockCopyFromContainer *mocks.MockDockerClient
	expected              string
	expectError           bool
}{
	{
		title:       "Success",
		containerID: "testid_1",
		mockCopyFromContainer: &mocks.MockDockerClient{
			CopyFromContainerFn: osReleaseCopyFromContainer("/etc/os-release", "DUMMYKEY1=dummyval1\nID=debian\nDUMMYKEY2=dummyval2"),
		},
		expected:    "debian",
		expectError: false,
	},
	{
		title:       "No ID found",
		containerID: "testid_2",
		mockCopyFromContainer: &mocks.MockDockerClient{
			CopyFromContainerFn: osReleaseCopyFromContainer("/etc/os-release", "DUMMYKEY1=dummyval1\nDUMMYKEY2=dummyval2"),
		},
		expected:    "unknown",
		expectError: false,
	},
	{
		title:       "newlines midfile in /etc/os-release",
		containerID: "testid_3",
		mockCopyFromContainer: &mocks.MockDockerClient{
			CopyFromContainerFn: osReleaseCopyFromContainer("/etc/os-release", "DUMMYKEY1=dummyval1\n\n\nID=alpine\nDUMMYKEY2=dummyval2"),
		},
		expected:    "alpine",
		expectError: false,
	},
	{
		title:       "newlines at beginning in /etc/os-release",
		containerID: "testid_4",
		mockCopyFromContainer: &mocks.MockDockerClient{
			CopyFromContainerFn: osReleaseCopyFromContainer("/etc/os-release", "\n\n\nID=fedora\nDUMMYKEY2=dummyval2"),
		},
		expected:    "fedora",
		expectError: false,
	},
	{
		title:       "newlines at end in /etc/os-release",
		containerID: "testid_5",
		mockCopyFromContainer: &mocks.MockDockerClient{
			CopyFromContainerFn: osReleaseCopyFromContainer("/etc/os-release", "ID=rhel\nDUMMYKEY2=dummyval2\n\n\n"),
		},
		expected:    "rhel",
		expectError: false,
	},
	{
		title:       "/etc/os-release not found, assuming scratch container",
		containerID: "testid_6",
		mockCopyFromContainer: &mocks.MockDockerClient{
			CopyFromContainerFn: fileNotFound(),
		},
		expected:    "scratch",
		expectError: false,
	},
	{
		title:       "error copying file from target, assuming unknown container",
		containerID: "testid_6",
		mockCopyFromContainer: &mocks.MockDockerClient{
			CopyFromContainerFn: errorCopyingFile(),
		},
		expected:    "unknown",
		expectError: true,
	},
	{
		title:       "ID key and value surrounded by different quotes with trailing and leading spaces",
		containerID: "testid",
		mockCopyFromContainer: &mocks.MockDockerClient{
			CopyFromContainerFn: osReleaseCopyFromContainer("/etc/os-release", "\"ID\"='   opensuse   '\nDUMMYKEY2=dummyval2\n\n\n"),
		},
		expected:    "opensuse",
		expectError: false,
	},
}

func Test_getContainerOS(t *testing.T) {
	for _, test := range getContainerOSTests {
		t.Run(test.title, func(t *testing.T) {
			actual, err := getContainerOS(context.TODO(), test.containerID, test.mockCopyFromContainer)
			if err != nil && !test.expectError {
				t.Errorf("Unexpected error returned: %s", err)
				t.FailNow()
			}
			if err == nil && test.expectError {
				t.Errorf("Expected error notreturned: %s", err)
				t.FailNow()
			}
			if actual != test.expected {
				t.Errorf("Test [%s] getContainerOS test failed, returned result [%s], does not match expected result [%s]",
					test.title,
					actual,
					test.expected,
				)
			}
		})
	}
}

// osReleaseCopyFromContainer returns a test closure CopyFromContainer function,
// which takes a string, s, to inject into the specified file, f, which is then
// returned in a tar archive within a ReadCloser
func osReleaseCopyFromContainer(f, s string) func(ctx context.Context, container, srcPath string, followSymlink bool) (io.ReadCloser, dockertypes.ContainerPathStat, error) {
	return func(ctx context.Context, container, srcPath string, followSymlink bool) (io.ReadCloser, dockertypes.ContainerPathStat, error) {
		containerPathStat := dockertypes.ContainerPathStat{
			LinkTarget: "",
		}
		return createTar(f, s), containerPathStat, nil
	}
}

// fileNotFound returns a test closure CopyFromContainer function, which
// simulates the response from CopyFromContainer when the requested file
// has not been found on the target container
func fileNotFound() func(ctx context.Context, container, srcPath string, followSymlink bool) (io.ReadCloser, dockertypes.ContainerPathStat, error) {
	return func(ctx context.Context, container, srcPath string, followSymlink bool) (io.ReadCloser, dockertypes.ContainerPathStat, error) {
		containerPathStat := dockertypes.ContainerPathStat{}
		err := fmt.Errorf("Simulated file not found error: no such file or directory")
		return nil, containerPathStat, err
	}
}

// errorCopyingFile returns a test closure CopyFromContainer function, which
// simulates a generic error response from CopyFromContainer
func errorCopyingFile() func(ctx context.Context, container, srcPath string, followSymlink bool) (io.ReadCloser, dockertypes.ContainerPathStat, error) {
	return func(ctx context.Context, container, srcPath string, followSymlink bool) (io.ReadCloser, dockertypes.ContainerPathStat, error) {
		containerPathStat := dockertypes.ContainerPathStat{}
		err := fmt.Errorf("Generic error response")
		return nil, containerPathStat, err
	}
}

// createTar returns a ReadCloser containing a tar archive
// with the specified file and its body
func createTar(filename, body string) io.ReadCloser {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	hdr := &tar.Header{
		Name: filename,
		Mode: 0600,
		Size: int64(len(body)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		log.Fatalln(err)
	}
	if _, err := tw.Write([]byte(body)); err != nil {
		log.Fatalln(err)
	}

	if err := tw.Close(); err != nil {
		log.Fatalln(err)
	}
	return ioutil.NopCloser(buf)
}
