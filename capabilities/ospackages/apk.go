package ospackages

import (
	"github.com/puppetlabs/lumogon/capabilities/payloadfilter"
	"github.com/puppetlabs/lumogon/capabilities/registry"
	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/types"
)

var apkDescription = `The apk capability returns all apk-managed package
names installed on a container, returning the results as a map:

map[string]string`

// The apk capability output from the container runtime exec
var apkCapability = dockeradapter.DockerAPICapability{
	Capability: types.Capability{
		Schema:      "http://puppet.com/lumogon/capability/label/draft-01/schema#1",
		Title:       "Packages (APK)",
		Name:        "apk",
		Description: apkDescription,
		Type:        "dockerapi",
		Payload:     nil,
		SupportedOS: map[string]int{"alpine": 1},
	},
	Harvest: func(capability *dockeradapter.DockerAPICapability, client dockeradapter.Harvester, id string, target types.TargetContainer) {
		logging.Stderr("[Apk] Harvesting packages from target %s [%s], harvester id: %s", target.Name, target.ID, id)
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
