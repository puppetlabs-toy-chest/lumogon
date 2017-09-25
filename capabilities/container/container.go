package container

import (
	"context"
	"fmt"
	"strings"

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
		logging.Debug("[Container Info] Harvesting container information associated with %s [%s]\n", target.Name, target.ID)
		capability.HarvestID = id

		containerInfo, err := InspectContainer(client, target.ID)
		if err != nil {
			capability.PayloadError(err.Error())
			return
		}

		filtered, _ := payloadfilter.Filter(containerInfo)

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

	var memorySwappiness int64
	if c.HostConfig.Resources.MemorySwappiness != nil {
		memorySwappiness = *c.HostConfig.Resources.MemorySwappiness
	}

	var oomKillDisable bool
	if c.HostConfig.Resources.OomKillDisable != nil {
		oomKillDisable = *c.HostConfig.Resources.OomKillDisable
	}

	result := map[string]interface{}{
		"hostname":           c.Config.Hostname,
		"domainname":         c.Config.Domainname,
		"user":               c.Config.User,
		"image":              c.Config.Image,
		"attachstdin":        fmt.Sprintf("%t", c.Config.AttachStdin),
		"attachstdout":       fmt.Sprintf("%t", c.Config.AttachStdout),
		"attachstderr":       fmt.Sprintf("%t", c.Config.AttachStderr),
		"tty":                fmt.Sprintf("%t", c.Config.Tty),
		"openstdin":          fmt.Sprintf("%t", c.Config.OpenStdin),
		"stdinonce":          fmt.Sprintf("%t", c.Config.StdinOnce),
		"privileged":         fmt.Sprintf("%t", c.HostConfig.Privileged),
		"publishallports":    fmt.Sprintf("%t", c.HostConfig.PublishAllPorts),
		"readonlyrootfs":     fmt.Sprintf("%t", c.HostConfig.ReadonlyRootfs),
		"shmsize":            c.HostConfig.ShmSize,
		"capadd":             strings.Join(c.HostConfig.CapAdd, ", "),
		"capdrop":            strings.Join(c.HostConfig.CapDrop, ", "),
		"runtime":            c.HostConfig.Runtime,
		"cpushares":          c.HostConfig.Resources.CPUShares,
		"memory":             c.HostConfig.Resources.Memory,
		"nanocpus":           c.HostConfig.Resources.NanoCPUs,
		"cpuperiod":          c.HostConfig.Resources.CPUPeriod,
		"cpuquota":           c.HostConfig.Resources.CPUQuota,
		"cpurealtimeperiod":  c.HostConfig.Resources.CPURealtimePeriod,
		"cpurealtimeruntime": c.HostConfig.Resources.CPURealtimeRuntime,
		"cpusetcpus":         c.HostConfig.Resources.CpusetCpus,
		"cpusetmems":         c.HostConfig.Resources.CpusetMems,
		"diskquota":          c.HostConfig.Resources.DiskQuota,
		"kernelmemory":       c.HostConfig.Resources.KernelMemory,
		"memoryreservation":  c.HostConfig.Resources.MemoryReservation,
		"memoryswap":         c.HostConfig.Resources.MemorySwap,
		"memoryswappiness":   memorySwappiness,
		"oomkilldisable":     fmt.Sprintf("%t", oomKillDisable),
		"pidslimit":          c.HostConfig.Resources.PidsLimit,
	}

	logging.Debug("[Container Info] Harvested [%+v]\n", result)

	return result, nil
}

func init() {
	logging.Debug("[Container Info] Initialising capability: %s\n", containerCapability.Title)
	registry.Registry.Add(containerCapability)
}
