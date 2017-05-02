package harvester

import (
	"time"

	"github.com/puppetlabs/transparent-containers/cli/types"
	"github.com/puppetlabs/transparent-containers/cli/utils"
)

// GenerateContainerReport gathers capabilities data for a specific container
func GenerateContainerReport(target types.TargetContainer, capabilityData map[string]types.Capability) *types.ContainerReport {
	containerReport := types.ContainerReport{
		Schema:            "http://puppet.com/lumogon/containerreport/draft-01/schema#1",
		Generated:         time.Now().String(),
		ContainerReportID: utils.GenerateUUID4(),
		ContainerID:       target.ID,
		ContainerName:     target.Name,
		Capabilities:      capabilityData,
	}
	return &containerReport
}
