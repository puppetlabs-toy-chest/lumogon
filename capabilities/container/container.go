package container

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/go-connections/nat"
	"github.com/puppetlabs/lumogon/capabilities/payloadfilter"
	"github.com/puppetlabs/lumogon/capabilities/registry"
	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/types"
)

var containerDescription = `The 'container' capability captures detailed container information`

var containerCapability = dockeradapter.DockerAPICapability{
	Capability: types.Capability{
		Schema:      "http://puppet.com/lumogon/capability/container/draft-01/schema#1",
		Title:       "Container Information",
		Name:        "container",
		Description: containerDescription,
		Type:        "dockerapi",
		Payload:     nil,
		SupportedOS: map[string]int{"all": 1},
	},
	Harvest: func(capability *dockeradapter.DockerAPICapability, client dockeradapter.Harvester, id string, target types.TargetContainer) {
		logging.Debug("[Container Info] Harvesting container information associated with %s [%s]", target.Name, target.ID)
		capability.HarvestID = id

		version, err := InspectContainer(client, target.ID)
		if err != nil {
			capability.PayloadError(err.Error())
			return
		}

		filtered, _ := payloadfilter.Filter(version)

		capability.Payload = filtered
	},
}

// InspectContainer Extracts and returns a formatted map[string]interface{} containing
// a subset of information returned by ContainerInspect
func InspectContainer(client dockeradapter.Harvester, targetID string) (map[string]interface{}, error) {
	ctx := context.Background()
	c, err := client.ContainerInspect(ctx, targetID)
	if err != nil {
		return nil, err
	}

	// TODO - **IMPORTANT** this contains only a subset of the information available
	// it explicitly avoids including any structured data (ports/mappings etc) pending
	// support in the UI, it also avoids any config that could potentially contain
	// sensitive data.

	result := map[string]interface{}{
		"Hostname":           c.Config.Hostname,
		"Domainname":         c.Config.Domainname,
		"User":               c.Config.User,
		"Image":              c.Config.Image,
		"AttachStdin":        fmt.Sprintf("%t", c.Config.AttachStdin),
		"AttachStdout":       fmt.Sprintf("%t", c.Config.AttachStdout),
		"AttachStderr":       fmt.Sprintf("%t", c.Config.AttachStderr),
		"Tty":                fmt.Sprintf("%t", c.Config.Tty),
		"OpenStdin":          fmt.Sprintf("%t", c.Config.OpenStdin),
		"StdinOnce":          fmt.Sprintf("%t", c.Config.StdinOnce),
		"Privileged":         fmt.Sprintf("%t", c.HostConfig.Privileged),
		"PublishAllPorts":    fmt.Sprintf("%t", c.HostConfig.PublishAllPorts),
		"ReadonlyRootfs":     fmt.Sprintf("%t", c.HostConfig.ReadonlyRootfs),
		"ShmSize":            fmt.Sprintf("%d", c.HostConfig.ShmSize),
		"CapAdd":             strings.Join(c.HostConfig.CapAdd, ", "),
		"CapDrop":            strings.Join(c.HostConfig.CapDrop, ", "),
		"Runtime":            c.HostConfig.Runtime,
		"CPUShares":          fmt.Sprintf("%d", c.HostConfig.Resources.CPUShares),
		"Memory":             fmt.Sprintf("%d", c.HostConfig.Resources.Memory),
		"NanoCPUs":           fmt.Sprintf("%d", c.HostConfig.Resources.NanoCPUs),
		"CPUPeriod":          fmt.Sprintf("%d", c.HostConfig.Resources.CPUPeriod),
		"CPUQuota":           fmt.Sprintf("%d", c.HostConfig.Resources.CPUQuota),
		"CPURealtimePeriod":  fmt.Sprintf("%d", c.HostConfig.Resources.CPURealtimePeriod),
		"CPURealtimeRuntime": fmt.Sprintf("%d", c.HostConfig.Resources.CPURealtimeRuntime),
		"CpusetCpus":         c.HostConfig.Resources.CpusetCpus,
		"CpusetMems":         c.HostConfig.Resources.CpusetMems,
		"DiskQuota":          fmt.Sprintf("%d", c.HostConfig.Resources.DiskQuota),
		"KernelMemory":       fmt.Sprintf("%d", c.HostConfig.Resources.KernelMemory),
		"MemoryReservation":  fmt.Sprintf("%d", c.HostConfig.Resources.MemoryReservation),
		"MemorySwap":         fmt.Sprintf("%d", c.HostConfig.Resources.MemorySwap),
		"MemorySwappiness":   fmt.Sprintf("%d", *c.HostConfig.Resources.MemorySwappiness),
		"OomKillDisable":     fmt.Sprintf("%t", *c.HostConfig.Resources.OomKillDisable),
		"PidsLimit":          fmt.Sprintf("%d", c.HostConfig.Resources.PidsLimit),
	}

	logging.Debug("[Container Info] Harvested [%+v]", result)

	return result, nil
}

func ports(m nat.PortSet) []string {
	keys := make([]string, len(m))

	i := 0
	for k := range m {
		keys[i] = k.Port()
		i++
	}
	return keys
}

func init() {
	logging.Debug("[Container Info] Initialising capability: %s", containerCapability.Title)
	registry.Registry.Add(containerCapability)
}
