package ospackages

import (
	"github.com/puppetlabs/lumogon/capabilities/payloadfilter"
	"github.com/puppetlabs/lumogon/capabilities/registry"
	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/types"
)

var dpkgDescription = `The dpkg capability returns all dpkg-managed package
names installed on a container, returning the results as a map:

map[string]string`

// The Dpkg capability output from the container runtime exec
var dpkgCapability = dockeradapter.DockerAPICapability{
	Capability: types.Capability{
		Schema:      "http://puppet.com/lumogon/capability/label/draft-01/schema#1",
		Title:       "Packages (DPKG)",
		Name:        "dpkg",
		Description: dpkgDescription,
		Type:        "dockerapi",
		Payload:     nil,
		SupportedOS: map[string]int{"ubuntu": 1, "debian": 1},
	},
	Harvest: func(capability *dockeradapter.DockerAPICapability, client dockeradapter.Harvester, id string, target types.TargetContainer) {
		logging.Debug("[Dpkg] Harvesting packages from target %s [%s], harvester id: %s", target.Name, target.ID, id)
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
	logging.Debug("[Dpkg] Initialising capability: %s", dpkgCapability.Title)
	registry.Registry.Add(dpkgCapability)
}
