package ospackages

import (
	"github.com/puppetlabs/transparent-containers/cli/capabilities/payloadfilter"
	"github.com/puppetlabs/transparent-containers/cli/capabilities/registry"
	"github.com/puppetlabs/transparent-containers/cli/dockeradapter"
	"github.com/puppetlabs/transparent-containers/cli/logging"
	"github.com/puppetlabs/transparent-containers/cli/types"
)

var dpkgDescription = `The dpkg capability returns all dpkg-managed package
names installed on a container, returning the results as a map:

map[string]string`

// The Dpkg capability output from the container runtime exec
var dpkgCapability = types.DockerAPICapability{
	Capability: types.Capability{
		Schema:      "http://puppet.com/lumogon/capability/label/draft-01/schema#1",
		Title:       "Dpkg Capability",
		Name:        "dpkg",
		Description: dpkgDescription,
		Type:        "dockerapi",
		Payload:     nil,
	},
	Harvest: func(capability *types.DockerAPICapability, client dockeradapter.Harvester, id string, target types.TargetContainer) {
		logging.Stderr("[Dpkg] Harvesting packages from %s [%s], harvestid: %s", id)
		capability.HarvestID = id

		output := make(map[string]interface{})
		result, err := runPackageCmd(client, target.ID, []string{"/bin/sh", "-c", `test -x /usr/bin/dpkg-query && dpkg-query -W -f '${Package},${Version}\n'`})
		if err != nil {
			capability.PayloadError(err.Error())
			return
		}

		for k, v := range result {
			output[k] = v
		}

		filtered, _ := payloadfilter.Filter(output)
		capability.Payload = filtered
	},
}

func init() {
	logging.Stderr("[Dpkg] Initialising capability: %s", dpkgCapability.Title)
	registry.Registry.Add(dpkgCapability)
}
