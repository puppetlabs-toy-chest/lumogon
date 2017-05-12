package scheduler

import (
	"testing"

	"fmt"

	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/test/registry"
	"github.com/puppetlabs/lumogon/types"
)

var getExpectedResultCountTests = []struct {
	title       string
	targets     []*types.TargetContainer
	registry    registry.MockRegistry
	expected    int
	expectError bool
}{
	{
		title: "AttachedCapabilities, all supportedOS, multiple targets",
		targets: []*types.TargetContainer{
			{
				ID:   "target1ID",
				Name: "target1Name",
				OSID: "debian",
			},
			{
				ID:   "target2ID",
				Name: "target2Name",
				OSID: "alpine",
			},
		},
		registry: registry.MockRegistry{
			AttachedCapabilitiesFn:  attachedCapabilities([]string{"all", "all"}),
			DockerAPICapabilitiesFn: dockerAPICapabilities([]string{}),
		},
		expected:    2,
		expectError: false,
	},
	{
		title: "DockerAPICapabilities, all supportedOS, multiple targets",
		targets: []*types.TargetContainer{
			{
				ID:   "target1ID",
				Name: "target1Name",
				OSID: "debian",
			},
			{
				ID:   "target2ID",
				Name: "target2Name",
				OSID: "alpine",
			},
		},
		registry: registry.MockRegistry{
			AttachedCapabilitiesFn:  attachedCapabilities([]string{}),
			DockerAPICapabilitiesFn: dockerAPICapabilities([]string{"all", "all"}),
		},
		expected:    2,
		expectError: false,
	},
	{
		title: "DockerAPICapabilities and AttachedCapabilities, all supportedOS, multiple targets",
		targets: []*types.TargetContainer{
			{
				ID:   "target1ID",
				Name: "target1Name",
				OSID: "debian",
			},
			{
				ID:   "target2ID",
				Name: "target2Name",
				OSID: "alpine",
			},
		},
		registry: registry.MockRegistry{
			AttachedCapabilitiesFn:  attachedCapabilities([]string{"all"}),
			DockerAPICapabilitiesFn: dockerAPICapabilities([]string{"all", "all"}),
		},
		expected:    4,
		expectError: false,
	},
	{
		title: "DockerAPICapabilities and AttachedCapabilities, only API with matching supported OS",
		targets: []*types.TargetContainer{
			{
				ID:   "target1ID",
				Name: "target1Name",
				OSID: "debian",
			},
			{
				ID:   "target2ID",
				Name: "target2Name",
				OSID: "alpine",
			},
		},
		registry: registry.MockRegistry{
			AttachedCapabilitiesFn:  attachedCapabilities([]string{"not_a_matching_os"}),
			DockerAPICapabilitiesFn: dockerAPICapabilities([]string{"debian", "alpine"}),
		},
		expected:    2,
		expectError: false,
	},
	{
		title: "AttachedCapabilities match subset of targets",
		targets: []*types.TargetContainer{
			{
				ID:   "target1ID",
				Name: "target1Name",
				OSID: "debian",
			},
			{
				ID:   "target2ID",
				Name: "target2Name",
				OSID: "alpine",
			},
		},
		registry: registry.MockRegistry{
			AttachedCapabilitiesFn:  attachedCapabilities([]string{"debian", "not_a_matching_os"}),
			DockerAPICapabilitiesFn: dockerAPICapabilities([]string{}),
		},
		expected:    1,
		expectError: false,
	},
	{
		title: "DockerAPICapabilities match subset of targets",
		targets: []*types.TargetContainer{
			{
				ID:   "target1ID",
				Name: "target1Name",
				OSID: "debian",
			},
			{
				ID:   "target2ID",
				Name: "target2Name",
				OSID: "alpine",
			},
		},
		registry: registry.MockRegistry{
			AttachedCapabilitiesFn:  attachedCapabilities([]string{}),
			DockerAPICapabilitiesFn: dockerAPICapabilities([]string{"alpine", "not_a_matching_os"}),
		},
		expected:    1,
		expectError: false,
	},
	{
		title: "No valid capabilities",
		targets: []*types.TargetContainer{
			{
				ID:   "target1ID",
				Name: "target1Name",
				OSID: "debian",
			},
			{
				ID:   "target2ID",
				Name: "target2Name",
				OSID: "alpine",
			},
		},
		registry: registry.MockRegistry{
			AttachedCapabilitiesFn:  attachedCapabilities([]string{"not_a_matching_os", "still_not_a_matching_os"}),
			DockerAPICapabilitiesFn: dockerAPICapabilities([]string{"not_a_matching_os", "still_not_a_matching_os"}),
		},
		expected:    0,
		expectError: false,
	},
}

// attachedCapabilities returns a test closure AttachedCapabilities function,
// which takes a slice of strings which are used to control the number of
// AttachedCapabilities returned and their SupportedOS
func attachedCapabilities(supported []string) func() []types.AttachedCapability {
	return func() []types.AttachedCapability {
		result := []types.AttachedCapability{}
		for i, os := range supported {
			c := types.AttachedCapability{
				Capability: types.Capability{
					Title:       fmt.Sprintf("AttachedCapability%d", i),
					SupportedOS: map[string]int{os: 1},
				},
			}
			result = append(result, c)
		}
		return result
	}
}

// dockerAPICapabilities returns a test closure DockerAPICapabilities function,
// which takes a slice of strings which are used to control the number of
// DockerAPICapabilities returned and their SupportedOS
func dockerAPICapabilities(supported []string) func() []dockeradapter.DockerAPICapability {
	return func() []dockeradapter.DockerAPICapability {
		result := []dockeradapter.DockerAPICapability{}
		for i, os := range supported {
			c := dockeradapter.DockerAPICapability{
				Capability: types.Capability{
					Title:       fmt.Sprintf("DockerAPICapability%d", i),
					SupportedOS: map[string]int{os: 1},
				},
			}
			result = append(result, c)
		}
		return result
	}
}

func Test_getExpectedResultCount(t *testing.T) {
	for _, test := range getExpectedResultCountTests {
		t.Run(test.title, func(t *testing.T) {
			actual := getExpectedResultCount(test.targets, test.registry)
			if actual != test.expected {
				t.Errorf("Test [%s] getExpectedResultCount test failed, returned result [%v], does not match expected result [%v]",
					test.title,
					actual,
					test.expected,
				)
			}
		})
	}
}
