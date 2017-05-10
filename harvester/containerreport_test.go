package harvester

import (
	"fmt"
	"testing"

	"github.com/puppetlabs/lumogon/types"
)

func TestGenerateContainerReport(t *testing.T) {
	// TODO this test is of limited value/terrible, delete?
	target := types.TargetContainer{
		ID:   "testid",
		Name: "testname",
	}
	capabiltyData := make(map[string]types.Capability)
	capabiltyData["capability_1"] = generateDummyCapability(1)
	capabiltyData["capability_2"] = generateDummyCapability(2)

	generatedReport := *GenerateContainerReport(target, capabiltyData)

	expectedSchema := "http://puppet.com/lumogon/containerreport/draft-01/schema#1"

	if generatedReport.Schema != expectedSchema {
		t.Errorf("Generated report schema [%s] does not match expected value [%s]", generatedReport.Schema, expectedSchema)
	}
	if generatedReport.ContainerID != target.ID {
		t.Errorf("Generated ContainerID [%s] does not match expected value [%s]", generatedReport.ContainerID, target.ID)
	}
	if generatedReport.ContainerName != target.Name {
		t.Errorf("Generated ContainerName [%s] does not match expected value [%s]", generatedReport.ContainerName, target.Name)
	}
	if len(generatedReport.Capabilities) != len(capabiltyData) {
		t.Errorf("Number of Capabilities [%d] does not match expected value [%d]", len(generatedReport.Capabilities), len(capabiltyData))
	}
}

func generateDummyCapability(id int) types.Capability {
	capability := types.Capability{
		Schema:      fmt.Sprintf("testschema_%d", id),
		Title:       fmt.Sprintf("testtitle_%d", id),
		Name:        fmt.Sprintf("testname_%d", id),
		Description: fmt.Sprintf("testdescription_%d", id),
		Type:        fmt.Sprintf("testtype_%d", id),
		HarvestID:   fmt.Sprintf("testharvestid_%d", id),
	}

	payload := make(map[string]interface{})
	payload[fmt.Sprintf("key_%d", id)] = []string{"value1", "value2"}

	capability.Payload = payload
	return capability
}
