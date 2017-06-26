package plugin

import (
	"github.com/johnmccabe/lumogon/dockeradapter"
	"github.com/johnmccabe/lumogon/types"
)

// Customtype TODO
type Customtype struct {
	StrF  string
	BoolF bool
	IntF  int
}

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

// LumogonPlugin TODO
type LumogonPlugin interface {
	Metadata() *Metadata
	Harvest(dockeradapter.Harvester, string, types.TargetContainer) (map[string]interface{}, error)
}
