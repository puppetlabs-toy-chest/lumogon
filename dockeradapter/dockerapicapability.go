package dockeradapter

import "github.com/puppetlabs/lumogon/types"

// DockerAPICapability embedded type adds a Docker specific Harvest function
// field which passes a client satisfying the dockeradapter.Harvester interface.
// This function is responsible for populating the Payload field.
type DockerAPICapability struct {
	types.Capability
	Harvest func(*DockerAPICapability, Harvester, string, types.TargetContainer) `json:"-"`
}
