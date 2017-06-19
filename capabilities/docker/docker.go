package docker

import (
	"context"

	"github.com/puppetlabs/lumogon/capabilities/payloadfilter"
	"github.com/puppetlabs/lumogon/capabilities/registry"
	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/types"
)

var dockerDescription = `The 'docker' capability captures information related to the underlying Docker server`

var dockerCapability = dockeradapter.DockerAPICapability{
	Capability: types.Capability{
		Schema:      "http://puppet.com/lumogon/capability/label/draft-01/schema#1",
		Title:       "Docker Server Information",
		Name:        "docker",
		Description: dockerDescription,
		Type:        "dockerapi",
		Payload:     nil,
		SupportedOS: map[string]int{"all": 1},
	},
	Harvest: func(capability *dockeradapter.DockerAPICapability, client dockeradapter.Harvester, id string, target types.TargetContainer) {
		logging.Debug("[Docker Server] Harvesting docker server information associated with %s [%s]", target.Name, target.ID)
		capability.HarvestID = id

		filtered, _ := payloadfilter.Filter(InfoToMap(client))

		capability.Payload = filtered
	},
}

// InfoToMap takes the Docker Context and returns a formatted map[string]interface{} contaning
// information exposted via the Docker API
func InfoToMap(client dockeradapter.Harvester) map[string]interface{} {
	ctx := context.Background()
	v := client.ServerVersion(ctx)

	return map[string]interface{}{
		"APIVersion":    v.APIVersion,
		"MinAPIVersion": v.MinAPIVersion,
	}
}

func init() {
	logging.Debug("[Docker Server] Initialising capability: %s", dockerCapability.Title)
	registry.Registry.Add(dockerCapability)
}
