package host

import (
	"fmt"

	"github.com/puppetlabs/lumogon/capabilities/payloadfilter"
	"github.com/puppetlabs/lumogon/capabilities/registry"
	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/types"
	gopsutilhost "github.com/shirou/gopsutil/host"
)

var hostDescription = `The 'host' capability captures the host.InfoStat from gopsutil
this includes:
- hostname
- uptime
- bootTime
- procs
- os
- platform
- platformFamily
- platformVersion
- kernelVersion
- virtualizationSystem
- virtualizationRole
- hostid`

var hostCapability = types.AttachedCapability{
	Capability: types.Capability{
		Schema:      "http://puppet.com/lumogon/capability/host/draft-01/schema#1",
		Title:       "Host Information",
		Name:        "host",
		Description: hostDescription,
		Type:        "attached",
		Payload:     nil,
		SupportedOS: map[string]int{"all": 1},
	},
	Harvest: func(capability *types.AttachedCapability, id string, args []string) {
		logging.Debug("[Host] Harvesting host capability, capability harvest id: %s", id)

		capability.HarvestID = id
		h, _ := gopsutilhost.Info()

		filtered, _ := payloadfilter.Filter(InfostatToMap(h))
		capability.Payload = filtered
	},
}

// InfostatToMap converts host.Infostat to map[string]interface{} for use in the ContainerReport
// sent to the harvester.
func InfostatToMap(h *gopsutilhost.InfoStat) map[string]interface{} {
	return map[string]interface{}{
		"hostname":             h.Hostname,
		"kernelversion":        h.KernelVersion,
		"os":                   h.OS,
		"procs":                fmt.Sprintf("%d", h.Procs),
		"platform":             h.Platform,
		"platformfamily":       h.PlatformFamily,
		"platformversion":      h.PlatformVersion,
		"uptime":               fmt.Sprintf("%d", h.Uptime),
		"virtualizationsystem": h.VirtualizationSystem,
		"virtualizationrole":   h.VirtualizationRole,
	}
}

func init() {
	logging.Debug("[Host] Initialising capability: %s", hostCapability.Title)
	registry.Registry.Add(hostCapability)
}
