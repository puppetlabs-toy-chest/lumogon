package registry

import (
	"context"
	"fmt"

	"os"

	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/types"
	"github.com/puppetlabs/lumogon/utils"
)

// CapabilitiesRegistry contains a record of each registered capability
// by type, Attached or Docker API
type CapabilitiesRegistry struct {
	attached  []types.AttachedCapability
	dockerAPI []dockeradapter.DockerAPICapability
}

// Registry provisioned capabilities
var Registry CapabilitiesRegistry

// Add registers a capability
func (c CapabilitiesRegistry) Add(capability interface{}) {
	switch capability.(type) {
	case types.AttachedCapability:
		logging.Stderr("[Registry] Adding ATTACHED capability to registry: %s\n", capability.(types.AttachedCapability).Title)
		Registry.attached = append(Registry.attached, capability.(types.AttachedCapability))
	case dockeradapter.DockerAPICapability:
		logging.Stderr("[Registry] Adding DOCKER API capability to registry: %s\n", capability.(dockeradapter.DockerAPICapability).Title)
		Registry.dockerAPI = append(Registry.dockerAPI, capability.(dockeradapter.DockerAPICapability))
	default:
		logging.Stdout("[Registry] Invalid capability type detected. Exiting..")
		os.Exit(1)
	}
}

// Count returns the total number of registered capabilities
func (c CapabilitiesRegistry) Count() int {
	return len(c.attached) + len(c.dockerAPI)
}

// TypesCount returns the number of registered types.
func (c CapabilitiesRegistry) TypesCount() int {
	var expectedTypes int
	if len(Registry.AttachedCapabilities()) > 0 {
		expectedTypes++
	}
	if len(Registry.DockerAPICapabilities()) > 0 {
		expectedTypes++
	}
	return expectedTypes
}

// AttachedCapabilities returns a list of AttachedCapability types
func (c CapabilitiesRegistry) AttachedCapabilities() []types.AttachedCapability {
	return c.attached
}

// DockerAPICapabilities returns a list of AttachedCapability types
func (c CapabilitiesRegistry) DockerAPICapabilities() []dockeradapter.DockerAPICapability {
	return c.dockerAPI
}

// DescribeCapability returns the description of a capability whose type is supplied
func (c CapabilitiesRegistry) DescribeCapability(capabilityID string) (string, error) {
	for _, attachedcapability := range Registry.AttachedCapabilities() {
		if attachedcapability.Name == capabilityID {
			return attachedcapability.Description, nil
		}
	}
	for _, dockerAPICapability := range Registry.DockerAPICapabilities() {
		if dockerAPICapability.Name == capabilityID {
			return dockerAPICapability.Description, nil
		}
	}

	return "", fmt.Errorf("[Registry] Unable to find capability: %s", capabilityID)
}

// Harvest gathers data registered capabilities with the specified harvest type
func Harvest(client dockeradapter.Harvester, targetContainerID string) map[string]types.Capability {
	harvestedData := map[string]types.Capability{}

	if client == nil {
		// Runs on the attached Harvester
		logging.Stderr("[Registry] Harvesting %d attached capabilities", len(Registry.AttachedCapabilities()))
		for _, attachedcapability := range Registry.AttachedCapabilities() {
			logging.Stderr("- %s\n", attachedcapability.Name)
			attachedcapability.Harvest(&attachedcapability, utils.GenerateUUID4(), []string{})
			harvestedData[attachedcapability.Name] = attachedcapability.Capability
		}
	}

	return harvestedData
}

// stringToTargetContainer converts a container ID or Name string into types.TargetContainer
func stringToTargetContainer(ctx context.Context, containerIDOrName string, client dockeradapter.Inspector) types.TargetContainer {
	containerJSON, err := client.ContainerInspect(ctx, containerIDOrName)
	if err != nil {
		error := fmt.Sprintf("[Registry] Unable to find target container: %s, error: %s", containerIDOrName, err)
		logging.Stderr(error)
	}
	targetContainer := types.TargetContainer{
		ID:   containerJSON.ContainerJSONBase.ID,
		Name: containerJSON.ContainerJSONBase.Name,
	}
	return targetContainer
}
