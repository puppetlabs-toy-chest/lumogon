package ospackages

import (
	"github.com/puppetlabs/lumogon/capabilities/payloadfilter"
	"github.com/puppetlabs/lumogon/capabilities/registry"
	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/types"
)

var rpmDescription = `The rpm capability returns all rpm-managed package
names installed on a container, returning the results as a map:

map[string]string`

// The rpm capability output from the container runtime exec
var rpmCapability = dockeradapter.DockerAPICapability{
	Capability: types.Capability{
		Schema:      "http://puppet.com/lumogon/capability/label/draft-01/schema#1",
		Title:       "Packages (RPM)",
		Name:        "rpm",
		Description: rpmDescription,
		Type:        "dockerapi",
		Payload:     nil,
		SupportedOS: map[string]int{"fedora": 1, "rhel": 1, "centos": 1, "opensuse": 1, "suse": 1, "ol": 1},
	},
	Harvest: func(capability *dockeradapter.DockerAPICapability, client dockeradapter.Harvester, id string, target types.TargetContainer) {
		logging.Debug("[Rpm] Harvesting packages from target %s [%s], harvester id: %s", target.Name, target.ID, id)
		capability.HarvestID = id

		output := make(map[string]interface{})
		result, err := runPackageCmd(client, target.ID, []string{"/bin/sh", "-c", `test -x /usr/bin/rpm && rpm  -qa --queryformat "%{NAME},%{VERSION}-%{RELEASE}-%{ARCH}\n"`})
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
	logging.Debug("[Rpm] Initialising capability: %s", rpmCapability.Title)
	registry.Registry.Add(rpmCapability)
}
