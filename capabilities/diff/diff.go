package diff

import (
	"context"
	"fmt"

	"github.com/docker/docker/pkg/archive"
	"github.com/puppetlabs/lumogon/capabilities/payloadfilter"
	"github.com/puppetlabs/lumogon/capabilities/registry"
	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/types"
)

var diffDescription = `The diff capability returns files changed from the initial image as a map["changed file"]"change type"`

// The Diff capability output from the container runtime inspect
var diffCapability = dockeradapter.DockerAPICapability{
	Capability: types.Capability{
		Schema:      "http://puppet.com/lumogon/capability/diff/draft-01/schema#1",
		Title:       "Changed Files",
		Name:        "diff",
		Description: diffDescription,
		Type:        "dockerapi",
		Payload:     nil,
		SupportedOS: map[string]int{"all": 1},
	},
	Harvest: func(capability *dockeradapter.DockerAPICapability, client dockeradapter.Harvester, id string, target types.TargetContainer) {
		logging.Debug("[Diff] Harvesting diff from %s [%s]", target.Name, target.ID)
		capability.HarvestID = id
		logging.Debug("[Diff] Harvesting diff capability, harvestid: %s", capability.HarvestID)

		ctx := context.Background()
		output := make(map[string]interface{})

		diffs, err := client.ContainerDiff(ctx, target.ID)
		if err != nil {
			errorMsg := fmt.Sprintf("[Diff] Error getting diff from targetContainer: %s, error: %s", target.Name, err)
			logging.Debug(errorMsg)
			capability.PayloadError(errorMsg)
			return
		}

		for _, diff := range diffs {
			logging.Debug("[Diff]   Path: %s, Kind %d", diff.Path, diff.Kind)
			var kind string
			switch diff.Kind {
			case archive.ChangeModify:
				kind = "Modified"
			case archive.ChangeAdd:
				kind = "Added"
			case archive.ChangeDelete:
				kind = "Deleted"
			}
			output[diff.Path] = kind
		}
		filtered, _ := payloadfilter.Filter(output)

		logging.Debug("[Diff]   Output: %v", output)
		logging.Debug("[Diff]   Filtered: %v", filtered)
		capability.Payload = filtered
	},
}

func init() {
	logging.Debug("[Diff] Initialising capability: %s", diffCapability.Title)
	registry.Registry.Add(diffCapability)
}
