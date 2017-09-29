package utils

import (
	"testing"
)

var getLocalContainerIDTests = []struct {
	title               string
	cgroupfile          string
	expectedContainerID string
	expectError         bool
}{
	{
		title:               "Simple success case",
		cgroupfile:          "./fixtures/cgroup.valid",
		expectedContainerID: "364b47437b435e097489faaafe9e896c20096fa555dae80d569e0cfa078a6d55",
		expectError:         false,
	},
	{
		title:               "cgroupCpuMultipleKeys",
		cgroupfile:          "./fixtures/cgroup.multi",
		expectedContainerID: "c7acb1eb27bb163d07c2ce0a3e7deb36a698824472f21074662e605bf57c7521",
		expectError:         false,
	},
	{
		title:               "nested container running in Docker-in-Docker",
		cgroupfile:          "./fixtures/cgroup.dind",
		expectedContainerID: "2ea749928409f83bd7af82f2d520f56345c99a4fcf5d0642cc99bd5cbc8168d5",
		expectError:         false,
	},
	{
		title:               "missingCgroupFile",
		cgroupfile:          "./fixtures/cgroup.doesnotexist",
		expectedContainerID: "n/a",
		expectError:         true,
	},
	{
		title:               "emptyCgroupFile",
		cgroupfile:          "./fixtures/cgroup.empty",
		expectedContainerID: "n/a",
		expectError:         true,
	},
	{
		title:               "cgroupCpuMissing",
		cgroupfile:          "./fixtures/cgroup.nocpu",
		expectedContainerID: "n/a",
		expectError:         true,
	},
	{
		title:               "cgroupCpuNoId",
		cgroupfile:          "./fixtures/cgroup.emptycpu",
		expectedContainerID: "n/a",
		expectError:         true,
	},
}

func Test_GetLocalContainerID(t *testing.T) {
	for _, test := range getLocalContainerIDTests {
		t.Run(test.title, func(t *testing.T) {
			actualContainerID, err := GetLocalContainerID(test.cgroupfile)
			if err != nil {
				if !test.expectError {
					t.Errorf("Test [%s] threw unexpected error [%s]", test.title, err)
				}
			} else {
				if actualContainerID != test.expectedContainerID {
					t.Errorf("Extracted containerID [%s] does not match expected value [%s]", actualContainerID, test.expectedContainerID)
				}
			}
		})
	}
}
