package capabilities

import (
	"github.com/puppetlabs/lumogon/capabilities/filesystem"
	"github.com/puppetlabs/lumogon/capabilities/host"
	"github.com/puppetlabs/lumogon/capabilities/label"
	"github.com/puppetlabs/lumogon/capabilities/ospackages"
)

// Init exists to allow capabilities init() functions to run when
// invoked from the Lumogon command handler.
func Init() {

	host.Init()
	label.Init()
	ospackages.Init()
	filesystem.Init()
}
