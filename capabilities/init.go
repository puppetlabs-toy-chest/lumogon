package capabilities

import (
	"github.com/puppetlabs/transparent-containers/cli/capabilities/host"
	"github.com/puppetlabs/transparent-containers/cli/capabilities/label"
	"github.com/puppetlabs/transparent-containers/cli/capabilities/ospackages"
)

// Init exists to allow capabilities init() functions to run when
// invoked from the Lumogon command handler.
func Init() {
	host.Init()
	label.Init()
	ospackages.Init()
}
