package label

import (
	"context"
	"fmt"

	"github.com/puppetlabs/transparent-containers/cli/capabilities/payloadfilter"
	"github.com/puppetlabs/transparent-containers/cli/capabilities/registry"
	"github.com/puppetlabs/transparent-containers/cli/dockeradapter"
	"github.com/puppetlabs/transparent-containers/cli/logging"
	"github.com/puppetlabs/transparent-containers/cli/types"
)

// Label contains a map of Docker labels
type Label struct {
	Labels map[string]string `json:"labels"`
}

var labelDescription = `The label capability returns all docker labels attached to a
container as a map:

map[string]string`

// The Label capability output from the container runtime inspect
var labelCapability = types.DockerAPICapability{
	Capability: types.Capability{
		Schema:      "http://puppet.com/lumogon/capability/label/draft-01/schema#1",
		Title:       "Label Capability",
		Name:        "label",
		Description: labelDescription,
		Type:        "dockerapi",
		Payload:     nil,
	},
	Harvest: func(capability *types.DockerAPICapability, client dockeradapter.Harvester, id string, target types.TargetContainer) {
		logging.Stderr("[Label] Harvesting label from %s [%s]", target.Name, target.ID)
		capability.HarvestID = id
		logging.Stderr("[Label]Harvesting label capability, harvestid: %s", capability.HarvestID)

		ctx := context.Background()
		output := make(map[string]interface{})

		containerJSON, err := client.ContainerInspect(ctx, target.ID)
		if err != nil {
			errorMsg := fmt.Sprintf("[Label] Error inspecting targetContainer: %s, error: %s", target.Name, err)
			logging.Stderr(errorMsg)
			capability.PayloadError(errorMsg)
			return
		}

		for k, v := range containerJSON.Config.Labels {
			output[k] = v
		}
		filtered, _ := payloadfilter.Filter(output)

		capability.Payload = filtered
	},
}

func init() {
	logging.Stderr("[Label] Initialising capability: %s", labelCapability.Title)
	registry.Registry.Add(labelCapability)
}
