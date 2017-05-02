package types

import "github.com/puppetlabs/lumogon/dockeradapter"

// Capability map of harvested capability data
type Capability struct {
	Schema      string                 `json:"$schema"`
	Title       string                 `json:"title"`
	Name        string                 `json:"-"`
	Description string                 `json:"-"`
	Type        string                 `json:"type"`
	HarvestID   string                 `json:"harvestid"`
	Payload     map[string]interface{} `json:"payload,omitempty"`
}

// DockerAPICapability embedded type adds a Docker specific Harvest function
// field which passes a client satisfying the dockeradapter.Harvester interface.
// This function is responsible for populating the Payload field.
type DockerAPICapability struct {
	Capability
	Harvest func(*DockerAPICapability, dockeradapter.Harvester, string, TargetContainer) `json:"-"`
}

// AttachedCapability embedded type adds a Harvest function field.
// This function is responsible for populating the Payload field.
type AttachedCapability struct {
	Capability
	Harvest func(*AttachedCapability, string, []string) `json:"-"`
}

// PayloadError records an error message in a capability.Payload
func (capability *Capability) PayloadError(message string) {
	capability.Payload = map[string]interface{}{"error": message}
}
