package harvester

import (
	"context"

	"sync"

	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/types"
	"github.com/puppetlabs/lumogon/utils"
)

// RunDockerAPIHarvester handles gathering DockerAPICapabilities from the target containers.
// It creates a ContainerReport for each target container populated with the output of the
// Harvest function from each capabilitie before sending the result to the collector via the
// main results channel, resultsCh.
func RunDockerAPIHarvester(ctx context.Context, wg *sync.WaitGroup, targets []*types.TargetContainer, capabilites []dockeradapter.DockerAPICapability, resultsCh chan types.ContainerReport, client dockeradapter.Harvester) error {
	defer logging.Stderr("[DockerAPI Harvester] Exiting")
	defer wg.Done()

	logging.Stderr("[DockerAPI Harvester] Running")
	if len(capabilites) == 0 {
		logging.Stderr("[DockerAPI Harvester] No Docker API Capabilities found")
		return nil
	}

	dockerAPIResultsCh := make(chan *types.ContainerReport)

	for _, target := range targets {
		go harvestDockerAPICapabilities(*target, client, capabilites, dockerAPIResultsCh)
	}

	for i := range targets {
		result := <-dockerAPIResultsCh
		logging.Stderr("[DockerAPI Harvester] Result [%d] received from name: %s, ID: %s", i, result.ContainerName, result.ContainerID)
		logging.Stderr("[DockerAPI Harvester] Sending to collector via resultsCh")
		resultsCh <- *result
	}

	return nil
}

func harvestDockerAPICapabilities(target types.TargetContainer, client dockeradapter.Harvester, capabilites []dockeradapter.DockerAPICapability, dockerAPIResultsCh chan *types.ContainerReport) {
	harvestedData := map[string]types.Capability{}

	logging.Stderr("[DockerAPI Harvester] Harvesting %d dockerAPI capabilities", len(capabilites))
	for _, capability := range capabilites {
		if !utils.KeyInMap("all", capability.SupportedOS) && !utils.KeyInMap(target.OSID, capability.SupportedOS) {
			logging.Stderr("[DockerAPI Harvester] skipping capability: %s, incompatible target OS: %s", capability.Name, target.OSID)
			continue
		}
		logging.Stderr("[DockerAPI Harvester] Harvesting %s\n", capability.Name)
		capability.Harvest(&capability, client, utils.GenerateUUID4(), target)
		logging.Stderr("[DockerAPI Harvester] Storing result %s\n", capability.Name)
		harvestedData[capability.Name] = capability.Capability
	}

	dockerAPIResultsCh <- GenerateContainerReport(target, harvestedData)
}
