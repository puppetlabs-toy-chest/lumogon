package utils

import (
	"strings"
	"testing"
)

func Test_GetLocalContainerID(t *testing.T) {
	cgroupfile := "./fixtures/cgroup.valid"
	expectedContainerID := "364b47437b435e097489faaafe9e896c20096fa555dae80d569e0cfa078a6d55"
	actualContainerID, err := GetLocalContainerID(cgroupfile)
	if err != nil {
		t.Errorf("Unexpected error returned [%s], expected containerID [%s]", err, expectedContainerID)
	}
	if strings.Compare(expectedContainerID, actualContainerID) != 0 {
		t.Errorf("Extracted containerID [%s] does not match expected value [%s]", actualContainerID, expectedContainerID)
	}
}

func Test_GetLocalContainerID_missingCgroupFile(t *testing.T) {
	cgroupfile := "./fixtures/cgroup.doesnotexist"
	_, err := GetLocalContainerID(cgroupfile)
	if err == nil {
		t.Errorf("Expected error not thrown")
	}
}

func Test_GetLocalContainerID_emptyCgroupFile(t *testing.T) {
	cgroupfile := "./fixtures/cgroup.empty"
	_, err := GetLocalContainerID(cgroupfile)
	if err == nil {
		t.Errorf("Expected error not thrown")
	}
}

func Test_GetLocalContainerID_cgroupCpuMissing(t *testing.T) {
	cgroupfile := "./fixtures/cgroup.nocpu"
	_, err := GetLocalContainerID(cgroupfile)
	if err == nil {
		t.Errorf("Expected error not thrown")
	}
}
