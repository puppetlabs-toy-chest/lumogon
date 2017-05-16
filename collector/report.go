package collector

import (
	"os"
	"time"

	"github.com/puppetlabs/lumogon/types"
	"github.com/puppetlabs/lumogon/version"
)

// NewReport gathers data from all capabilities
func NewReport(reportID string, clientVersion version.ClientVersion) *types.Report {
	hostname, _ := os.Hostname()
	report := types.Report{
		Schema:        "http://puppet.com/lumogon/core/draft-01/schema#1",
		Generated:     time.Now().String(),
		Host:          hostname,
		Owner:         "default",
		Group:         []string{"default"},
		ClientVersion: clientVersion,
		ReportID:      reportID,
		Containers:    map[string]types.ContainerReport{},
	}

	return &report
}
