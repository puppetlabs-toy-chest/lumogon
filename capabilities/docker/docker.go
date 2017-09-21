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
		Schema:      "http://puppet.com/lumogon/capability/docker/draft-01/schema#1",
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

		version, err := VersionToMap(client)
		if err != nil {
			capability.PayloadError(err.Error())
			return
		}

		filtered, _ := payloadfilter.Filter(version)

		capability.Payload = filtered
	},
}

// VersionToMap Extracts and returns a formatted map[string]interface{} containing
// version information exposted via the Docker API
func VersionToMap(client dockeradapter.Harvester) (map[string]interface{}, error) {
	ctx := context.Background()
	v, err := client.ServerVersion(ctx)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Version":       v.Version,
		"APIVersion":    v.APIVersion,
		"MinAPIVersion": v.MinAPIVersion,
	}, nil
}

func init() {
	logging.Debug("[Docker Server] Initialising capability: %s", dockerCapability.Title)
	registry.Registry.Add(dockerCapability)
}
