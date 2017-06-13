package filesystem

import (
	"context"
	"fmt"

	"github.com/puppetlabs/lumogon/capabilities/payloadfilter"
	"github.com/puppetlabs/lumogon/capabilities/registry"
	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/types"
)

var filesystemDescription = `The filesystem capability returns sizes of filesystem in a container as a map["layer"]"size"`

// The filesystemCapability capability output from the container runtime inspect
var filesystemCapability = dockeradapter.DockerAPICapability{
	Capability: types.Capability{
		Schema:      "http://puppet.com/lumogon/capability/diff/draft-01/schema#1",
		Title:       "Filesystem",
		Name:        "filesystem",
		Description: filesystemDescription,
		Type:        "dockerapi",
		Payload:     nil,
		SupportedOS: map[string]int{"all": 1},
	},
	Harvest: func(capability *dockeradapter.DockerAPICapability, client dockeradapter.Harvester, id string, target types.TargetContainer) {
		logging.Stderr("[Filesystem] Harvesting filesystem sizes from %s [%s]", target.Name, target.ID)
		capability.HarvestID = id
		logging.Stderr("[Filesystem] Harvesting filesystem capability, harvestid: %s", capability.HarvestID)

		ctx := context.Background()
		output := make(map[string]interface{})

		containerData, err := client.ContainerFilesystem(ctx, target.ID)
		if err != nil {
			errorMsg := fmt.Sprintf("[Filesystem] Error getting filesystem data from targetContainer: %s, error: %s", target.Name, err)
			logging.Stderr(errorMsg)
			capability.PayloadError(errorMsg)
			return
		}

		output["sizerw"] = containerData.SizeRw
		output["sizerootfs"] = containerData.SizeRootFs
		filtered, _ := payloadfilter.Filter(output)
		logging.Stderr("[Filesystem]   Output: %v", output)
		logging.Stderr("[Filesystem]   Filtered: %v", filtered)
		capability.Payload = filtered
	},
}

func init() {
	logging.Stderr("[Filesystem] Initialising capability: %s", filesystemCapability.Title)
	registry.Registry.Add(filesystemCapability)
}
