package version

import (
	"testing"
)

func TestBuildVersionFromMakefile(t *testing.T) {
	expectedBuildVersion := "testversionstring"
	if Version.BuildVersion != expectedBuildVersion {
		t.Errorf("BuildVersion [%s] does not match expected value [%s]", Version.BuildVersion, expectedBuildVersion)
	}
}

func TestBuildTimeFromMakefile(t *testing.T) {
	expectedBuildTime := "testdatestring"
	if Version.BuildTime != expectedBuildTime {
		t.Errorf("BuildTime [%s] does not match expected value [%s]", Version.BuildTime, expectedBuildTime)
	}
}

func TestBuildSHAFromMakefile(t *testing.T) {
	expectedBuildSHA := "testbuildsha"
	if Version.BuildSHA != expectedBuildSHA {
		t.Errorf("BuildSHA [%s] does not match expected value [%s]", Version.BuildSHA, expectedBuildSHA)
	}
}
