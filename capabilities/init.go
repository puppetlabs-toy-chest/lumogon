package capabilities

import (
	"github.com/puppetlabs/lumogon/capabilities/container"
	"github.com/puppetlabs/lumogon/capabilities/diff"
	"github.com/puppetlabs/lumogon/capabilities/docker"
	"github.com/puppetlabs/lumogon/capabilities/host"
	"github.com/puppetlabs/lumogon/capabilities/label"
	"github.com/puppetlabs/lumogon/capabilities/ospackages"
)

// Init exists to allow capabilities init() functions to run when
// invoked from the Lumogon command handler.
func Init() {
	docker.Init()
	host.Init()
	label.Init()
	ospackages.Init()
	diff.Init()
	container.Init()
}
