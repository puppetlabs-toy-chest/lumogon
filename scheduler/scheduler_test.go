package scheduler

import (
	"context"
	"testing"

	"fmt"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/puppetlabs/lumogon/test/mocks"
	"github.com/puppetlabs/lumogon/types"
)

func Test_stringToTargetContainer(t *testing.T) {
	var mockClient = mocks.MockDockerClient{
		ContainerInspectFn: func(ctx context.Context, containerID string) (dockertypes.ContainerJSON, error) {
			containerJSON := dockertypes.ContainerJSON{}
			containerJSON.ContainerJSONBase = &dockertypes.ContainerJSONBase{}
			containerJSON.ContainerJSONBase.ID = fmt.Sprintf("testID_%s", containerID)
			containerJSON.ContainerJSONBase.Name = fmt.Sprintf("testName_%s", containerID)
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

func Test_stringsToTargetContainers(t *testing.T) {
	var mockClient = mocks.MockDockerClient{
		ContainerInspectFn: func(ctx context.Context, containerID string) (dockertypes.ContainerJSON, error) {
			containerJSON := dockertypes.ContainerJSON{}
			containerJSON.ContainerJSONBase = &dockertypes.ContainerJSONBase{}
			containerJSON.ContainerJSONBase.ID = fmt.Sprintf("testID_%s", containerID)
			containerJSON.ContainerJSONBase.Name = fmt.Sprintf("testName_%s", containerID)
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
