package plugin

import (
	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/types"
)

// Types of Plugin
const (
	Attached int = iota
	DockerAPI
)

// Metadata TODO
type Metadata struct {
	Schema      string
	ID          string
	Name        string
	Description string
	Type        int
	Version     string
	GitSHA      string
	SupportedOS map[string]int
}

// Plugin TODO
type Plugin interface {
	Metadata() *Metadata
	Harvest(dockeradapter.Harvester, string, types.TargetContainer) (map[string]interface{}, error)
}
