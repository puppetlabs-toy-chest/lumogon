package ospackages

import (
	"github.com/puppetlabs/transparent-containers/cli/capabilities/payloadfilter"
	"github.com/puppetlabs/transparent-containers/cli/capabilities/registry"
	"github.com/puppetlabs/transparent-containers/cli/dockeradapter"
	"github.com/puppetlabs/transparent-containers/cli/logging"
	"github.com/puppetlabs/transparent-containers/cli/types"
)

var apkDescription = `The apk capability returns all apk-managed package
names installed on a container, returning the results as a map:

map[string]string`

// The apk capability output from the container runtime exec
var apkCapability = types.DockerAPICapability{
	Capability: types.Capability{
		Schema:      "http://puppet.com/lumogon/capability/label/draft-01/schema#1",
		Title:       "Apk Capability",
		Name:        "apk",
		Description: apkDescription,
		Type:        "dockerapi",
		Payload:     nil,
	},
	Harvest: func(capability *types.DockerAPICapability, client dockeradapter.Harvester, id string, target types.TargetContainer) {
		logging.Stderr("[Apk] Harvesting packages from %s [%s], harvestid: %s", id)
		capability.HarvestID = id

		output := make(map[string]interface{})
		cmd := []string{"/bin/sh", "-c", `test -x /sbin/apk && apk info -v | grep -v WARNING | sed  -r 's/(.*)-([0-9._azAZ]+)/\1,\2/'`}
		result, err := runPackageCmd(client, target.ID, cmd)
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
	logging.Stderr("[Apk] Initialising capability: %s", apkCapability.Title)
	registry.Registry.Add(apkCapability)
}
