package collector

import (
	"time"

	"github.com/puppetlabs/transparent-containers/cli/types"
	"github.com/puppetlabs/transparent-containers/cli/version"
)

// NewReport gathers data from all capabilities
func NewReport(reportID string, clientVersion version.ClientVersion) *types.Report {
	report := types.Report{
		Schema:        "http://puppet.com/lumogon/core/draft-01/schema#1",
		Generated:     time.Now().String(),
		Owner:         "default",
		Group:         []string{"default"},
		ClientVersion: clientVersion,
		ReportID:      reportID,
		Containers:    map[string]types.ContainerReport{},
	}

	return &report
}
